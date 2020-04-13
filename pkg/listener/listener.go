// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package listener

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	fabcfg "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/cryptoutil"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	"github.com/kataras/golog"
	"github.com/spf13/viper"

	"github.com/key-inside/patrasche/pkg/config"
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
		return config.PutStringWithARN(*keep.arn, keep.stringNumber())
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

	var cfg core.ConfigProvider
	arn, path, err := config.GetARN("network")
	if err != nil {
		cfg = fabcfg.FromFile(path)
	} else { // AWS resource
		in, typ, err := config.GetReaderWithARN(arn)
		if err != nil {
			return err
		}
		cfg = fabcfg.FromReader(in, typ)
	}

	// sdk
	sdk, err := fabsdk.New(cfg)
	if err != nil {
		return err
	}
	defer sdk.Close()

	// msp
	si, err := getSigningIdentity(sdk.Context())
	if err != nil {
		return err
	}

	// block number keepping
	keep, err := loadBlockKeep()
	if err != nil {
		return err
	}

	// channel provider
	ctx := sdk.ChannelContext(viper.GetString("channel"), fabsdk.WithIdentity(si))

	// options
	var opts []event.ClientOption
	if keep.newest {
		opts = []event.ClientOption{event.WithBlockEvents(), event.WithSeekType(seek.Newest)}
		// seek.Newest starting from lastest block.
		// so, you must have a plan to prevent duplicated processing the block
	} else {
		opts = []event.ClientOption{event.WithBlockEvents(), event.WithSeekType(seek.FromBlock), event.WithBlockNum(keep.number)}
	}

	// client & listen
	client, err := event.New(ctx, opts...)
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

	idFSet := viper.IsSet("tx.id")
	if idFSet {
		idFilter = regexp.MustCompilePOSIX(viper.GetString("tx.id"))
		golog.Infof("TX-Filter[ID] %s", idFilter.String())
	}
	typeFSet := viper.IsSet("tx.type")
	if typeFSet {
		typsStr := viper.GetString("tx.type")
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

	follow := viper.GetBool("follow")

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

func getSigningIdentity(ctx context.ClientProvider) (mspctx.SigningIdentity, error) {
	client, err := mspclient.New(ctx)
	if err != nil {
		return nil, err
	}

	arn, nameOrPath, err := config.GetARN("identity")
	if err != nil {
		var err error
		si, err := client.GetSigningIdentity(nameOrPath)
		if err != nil {
			if err != mspclient.ErrUserNotFound {
				golog.Debugf("%#v", err)
				return nil, err
			}

			// msp style
			golog.Debug("identity from msp directory")

			mspDir := nameOrPath
			_ctx, err := ctx()
			if err != nil {
				return nil, err
			}
			cert, err := ioutil.ReadFile(filepath.Join(mspDir, "signcerts", "cert.pem"))
			if err != nil {
				return nil, err
			}
			pubKey, err := cryptoutil.GetPublicKeyFromCert(cert, _ctx.CryptoSuite())
			if err != nil {
				return nil, err
			}
			priKey, _ := ioutil.ReadFile(filepath.Join(mspDir, "keystore", fmt.Sprintf("%x_sk", pubKey.SKI())))

			return client.CreateSigningIdentity(mspctx.WithCert(cert), mspctx.WithPrivateKey(priKey))
		}
		// else
		golog.Debug("identity from store")
		return si, nil
	}

	// else AWS resource
	golog.Debug("identity from arn")

	in, typ, err := config.GetReaderWithARN(arn)
	if err != nil {
		return nil, err
	}
	// uses viper for config consistency
	temp := viper.New()
	temp.SetConfigType(typ)
	if err := temp.ReadConfig(in); err != nil {
		return nil, err
	}
	cert := []byte(temp.GetString("cert"))
	priKey := []byte(temp.GetString("key"))

	return client.CreateSigningIdentity(mspctx.WithCert(cert), mspctx.WithPrivateKey(priKey))
}

func loadBlockKeep() (blockKeep, error) {
	keep := blockKeep{newest: true}
	keepNum := ""
	keepArn, numOrPath, err := config.GetARN("block")
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
		keepNum, err = config.GetStringWithARN(keepArn)
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
