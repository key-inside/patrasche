// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package listener

import (
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
	"github.com/kataras/golog"
	"github.com/spf13/viper"

	"github.com/key-inside/patrasche/pkg/aws"
	"github.com/key-inside/patrasche/pkg/channel"
	"github.com/key-inside/patrasche/pkg/channel/event"
	"github.com/key-inside/patrasche/pkg/proto"
	"github.com/key-inside/patrasche/pkg/tx"
)

type blockKeep struct {
	newest bool
	number uint64
	arn    *arn.ARN
	path   string
}

func (keep blockKeep) stringNumber() string {
	return strconv.FormatUint(keep.number, 10)
}

func (keep blockKeep) save() error {
	if keep.arn != nil {
		return aws.PutStringWithARN(*keep.arn, keep.stringNumber())
	} else if keep.path != "" {
		return ioutil.WriteFile(keep.path, []byte(keep.stringNumber()), 0644)
	}
	return nil
}

// Listen _
func Listen(txh tx.Handler) error {
	if nil == txh {
		golog.Fatal("nil handler")
	}

	// patrasche fabric channel
	channel, err := channel.New()
	if err != nil {
		return err
	}
	defer channel.Close()

	// block number keepping
	keep, err := loadBlockKeep()
	if err != nil {
		return err
	}

	// options
	// var opts []event.ClientOption
	var seekType seek.Type = seek.FromBlock
	if keep.newest {
		seekType = seek.Newest
	}

	// client & listen
	client, err := channel.NewEventClient(keep.number, seekType)
	if err != nil {
		return err
	}
	return listenBlockEvent(client, txh, txFilter(), keep)
}

// TxFilter _
type TxFilter func(*tx.Tx) bool

func txFilter() TxFilter {
	var idFilter *regexp.Regexp
	var typeFilter common.HeaderType

	idFSet := viper.IsSet("patrasche.tx.id")
	if idFSet {
		idFilter = regexp.MustCompilePOSIX(viper.GetString("patrasche.tx.id"))
		golog.Infof("TX-Filter[ID] %s", idFilter.String())
	}
	typeFSet := viper.IsSet("patrasche.tx.type")
	if typeFSet {
		typsStr := viper.GetString("patrasche.tx.type")
		if i, err := strconv.Atoi(typsStr); err == nil {
			typeFilter = common.HeaderType(i)
		} else {
			typeFilter = common.HeaderType(common.HeaderType_value[typsStr])
		}
		golog.Infof("TX-Filter[TYPE] %s", typeFilter.String())
	}

	return func(t *tx.Tx) bool {
		typeOk := !typeFSet || (t.HeaderType() == typeFilter)
		idOk := !idFSet || idFilter.MatchString(t.ID())
		return typeOk && idOk
	}
}

func listenBlockEvent(client *event.Client, txh tx.Handler, txFilter TxFilter, keep blockKeep) error {
	registration, notifier, err := client.RegisterBlockEvent()
	if err != nil {
		return err
	}
	defer client.Unregister(registration)

	follow := viper.GetBool("patrasche.follow")

	for {
		select {
		case e := <-notifier:
			block := e.Block
			blockNum := block.Header.Number
			blockHash, err := proto.GenerateBlockHash(block)
			if err != nil {
				return err
			}
			golog.Infof("BLOCK[%d] TxCount: %d Hash: %x", blockNum, len(block.Data.Data), blockHash)

			// keep current block number if needed
			keep.number = blockNum
			if err := keep.save(); err != nil {
				return err
			}

			for i, data := range block.Data.Data {
				header, transaction, err := unmarshalTx(data)
				if err != nil {
					return err
				}
				t := &tx.Tx{
					BlockNum:       blockNum,
					Seq:            i,
					Header:         header,
					Transaction:    transaction,
					ValidationCode: peer.TxValidationCode(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER][i]), // BlockMetadataIndex_TRANSACTIONS_FILTER = 2
				}
				typ := t.HeaderType()
				utc := t.UTC().Format("2006-01-02T15:04:05.000000000Z07:00")
				pass := txFilter(t)
				golog.Debugf("TX[%s] Type: %s(%d), ValidationCode: %s(%d), Timestamp: %s, Pass: %t", t.ID(), typ.String(), typ, t.ValidationCode.String(), t.ValidationCode, utc, pass)

				if pass {
					if err := txh.Handle(t); err != nil {
						return err
					}
				}
			}
		}
		if !follow {
			break
		}
	}

	return nil
}

func loadBlockKeep() (blockKeep, error) {
	keep := blockKeep{newest: true}
	keepNum := ""
	keepArn, numOrPath, err := aws.GetARN("patrasche.block")
	if err != nil {
		if numOrPath != "" {
			num, err := strconv.ParseUint(numOrPath, 10, 64) // check number or path
			if err != nil {                                  // path
				keep.path = numOrPath // keep file
				nBytes, err := ioutil.ReadFile(keep.path)
				if err != nil {
					if !os.IsNotExist(err) {
						return keep, err
					} // else ignore
				}
				keepNum = string(nBytes)
			} else { // start number without keep
				keep.number = num
				keep.newest = false
			}
		} // else nothing for block number
	} else { // AWS resource
		keep.arn = &keepArn // keep ARN
		keepNum, err = aws.GetStringWithARN(keepArn)
		if err != nil {
			return keep, err
		}
	}

	if keepNum != "" { // number from file or AWS
		keepNum = strings.TrimSpace(keepNum)
		num, err := strconv.ParseUint(keepNum, 10, 64)
		if err != nil {
			return keep, err
		}
		keep.number = num
		keep.newest = false
	}
	return keep, nil
}

func unmarshalTx(data []byte) (*common.ChannelHeader, *peer.Transaction, error) {
	envelope, err := proto.UnmarshalEnvelope(data)
	if err != nil {
		return nil, nil, err
	}
	payload, err := proto.UnmarshalPayload(envelope.Payload)
	if err != nil {
		return nil, nil, err
	}
	channelHeader, err := proto.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, nil, err
	}
	transaction, err := proto.UnmarshalTransaction(payload.Data)
	if err != nil {
		return nil, nil, err
	}
	return channelHeader, transaction, nil
}
