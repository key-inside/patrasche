// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package tx

import (
	"time"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/key-inside/patrasche/pkg/proto"
)

// Tx _
type Tx struct {
	BlockNum       uint64
	Seq            int
	Header         *common.ChannelHeader
	Transaction    *peer.Transaction
	ValidationCode peer.TxValidationCode
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
