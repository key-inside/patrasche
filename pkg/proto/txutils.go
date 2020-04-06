// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package proto

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
)

// GetPayloads _
func GetPayloads(txActions *peer.TransactionAction) (*peer.ChaincodeActionPayload, *peer.ChaincodeAction, error) {
	// TODO: pass in the tx type (in what follows we're assuming the
	// type is ENDORSER_TRANSACTION)
	ccPayload, err := UnmarshalChaincodeActionPayload(txActions.Payload)
	if err != nil {
		return nil, nil, err
	}

	if ccPayload.Action == nil || ccPayload.Action.ProposalResponsePayload == nil {
		return nil, nil, errors.New("no payload in ChaincodeActionPayload")
	}
	pRespPayload, err := UnmarshalProposalResponsePayload(ccPayload.Action.ProposalResponsePayload)
	if err != nil {
		return nil, nil, err
	}

	if pRespPayload.Extension == nil {
		return nil, nil, errors.New("response payload is missing extension")
	}

	respPayload, err := UnmarshalChaincodeAction(pRespPayload.Extension)
	if err != nil {
		return ccPayload, nil, err
	}
	return ccPayload, respPayload, nil
}

// GetTxReadWriteSet _
func GetTxReadWriteSet(protoBytes []byte) (*rwset.TxReadWriteSet, error) {
	rwSet := &rwset.TxReadWriteSet{}
	if err := proto.Unmarshal(protoBytes, rwSet); err != nil {
		return nil, err
	}
	return rwSet, nil
}

// GetKVRWSet _
func GetKVRWSet(protoBytes []byte) (*kvrwset.KVRWSet, error) {
	kvSet := &kvrwset.KVRWSet{}
	if err := proto.Unmarshal(protoBytes, kvSet); err != nil {
		return nil, err
	}
	return kvSet, nil
}
