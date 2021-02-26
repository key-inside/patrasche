// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package tx

import (
	"time"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/key-inside/patrasche/pkg/proto"
)

// Tx _
type Tx struct {
	BlockNum        uint64
	Seq             int
	Header          *common.ChannelHeader
	SignatureHeader *common.SignatureHeader
	Transaction     *peer.Transaction
	ValidationCode  peer.TxValidationCode
}

// New _
func New(blockNum uint64, seq int, validationByte byte, data []byte) (*Tx, error) {
	envelope, err := proto.UnmarshalEnvelope(data)
	if err != nil {
		return nil, err
	}
	payload, err := proto.UnmarshalPayload(envelope.Payload)
	if err != nil {
		return nil, err
	}
	channelHeader, err := proto.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, err
	}
	signatureHeader, err := proto.UnmarshalSignatureHeader(payload.Header.SignatureHeader)
	if err != nil {
		return nil, err
	}
	transaction, err := proto.UnmarshalTransaction(payload.Data)
	if err != nil {
		return nil, err
	}
	return &Tx{
		BlockNum:        blockNum,
		Seq:             seq,
		Header:          channelHeader,
		SignatureHeader: signatureHeader,
		Transaction:     transaction,
		ValidationCode:  peer.TxValidationCode(validationByte),
	}, nil
}

// ID returns TxID
func (t Tx) ID() string {
	return t.Header.TxId
}

// HeaderType _
func (t Tx) HeaderType() common.HeaderType {
	return common.HeaderType(t.Header.Type)
}

// IsValid _
func (t Tx) IsValid() bool {
	return peer.TxValidationCode_VALID == t.ValidationCode
}

// UTC _
func (t Tx) UTC() time.Time {
	ts := t.Header.Timestamp
	utc := time.Unix(ts.Seconds, int64(ts.Nanos)).UTC()
	return utc
}

// MSPID _
func (t Tx) MSPID() string {
	sid, err := t.GetIdentity()
	if err != nil {
		return ""
	}
	return sid.Mspid
}

// GetIdentity _
func (t Tx) GetIdentity() (*msp.SerializedIdentity, error) {
	sid, err := proto.UnmarshalSerializedIdentity(t.SignatureHeader.Creator)
	if err != nil {
		return nil, err
	}
	return sid, nil
}

// GetChaincodeAction _
func (t Tx) GetChaincodeAction() (*peer.ChaincodeAction, error) {
	if t.Transaction != nil && len(t.Transaction.Actions) > 0 {
		_, ccA, err := proto.GetPayloads(t.Transaction.Actions[0])
		if err != nil {
			return nil, err
		}
		return ccA, nil
	}
	return nil, nil
}

// GetChaincodeEvent _
func (t Tx) GetChaincodeEvent() (*peer.ChaincodeEvent, error) {
	ccA, err := t.GetChaincodeAction()
	if err != nil {
		return nil, err
	}
	if ccA != nil && ccA.Events != nil {
		return proto.UnmarshalChaincodeEvents(ccA.Events)
	}
	return nil, nil
}

// GetChaincodeInvocationSpec _
func (t Tx) GetChaincodeInvocationSpec() (*peer.ChaincodeInvocationSpec, error) {
	if t.Transaction != nil && len(t.Transaction.Actions) > 0 {
		ccAP, _, err := proto.GetPayloads(t.Transaction.Actions[0])
		if err != nil {
			return nil, err
		}
		ccP, err := proto.UnmarshalChaincodeProposalPayload(ccAP.ChaincodeProposalPayload)
		if err != nil {
			return nil, err
		}
		return proto.UnmarshalChaincodeInvocationSpec(ccP.Input)
	}
	return nil, nil
}

// GetReadWriteSet _
func (t Tx) GetReadWriteSet() (*rwset.TxReadWriteSet, error) {
	ccA, err := t.GetChaincodeAction()
	if err != nil {
		return nil, err
	}
	if ccA != nil {
		return proto.GetTxReadWriteSet(ccA.Results)
	}
	return nil, nil
}
