// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package inspect

import (
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/kataras/golog"

	"github.com/key-inside/patrasche/pkg/proto"
	"github.com/key-inside/patrasche/pkg/tx"
)

type handler struct{}

// Handle _
func (h handler) Handle(t *tx.Tx) error {
	if t.IsValid() && t.HeaderType() == common.HeaderType_ENDORSER_TRANSACTION {
		golog.Infof("MSPID: %s", t.MSPID())

		// ChaincodeInvocationSpec
		spec, err := t.GetChaincodeInvocationSpec()
		if err != nil {
			return err
		}
		golog.Infof("ChaincodeInvocationSpec.ChaincodeId: %#v", spec.ChaincodeSpec.ChaincodeId)
		for ai, arg := range spec.ChaincodeSpec.Input.Args {
			golog.Infof("ChaincodeInvocationSpec.Input.Args[%d]: %s", ai, string(arg))
		}

		// ChaincodeAction
		ccA, err := t.GetChaincodeAction()
		if err != nil {
			return err
		}
		golog.Infof("ChaincodeAction.ChaincodeId: %#v", ccA.ChaincodeId)
		golog.Infof("ChaincodeAction.Response: %#v", ccA.Response)
		golog.Infof("ChaincodeAction.Response.Payload: %s", string(ccA.Response.Payload))

		// ChaincodeEvent
		ccE, err := t.GetChaincodeEvent()
		if err != nil {
			return err
		}
		if ccE != nil {
			golog.Infof("ChaincodeAction.Events: %#v", ccE)
			golog.Infof("ChaincodeAction.Events.Payload: %s", string(ccE.Payload))
		}

		// RWSet
		rws, err := t.GetReadWriteSet()
		for _, nss := range rws.NsRwset {
			kvs, err := proto.GetKVRWSet(nss.Rwset)
			if err != nil {
				return err
			}
			golog.Infof("TxRwSet[%s]", nss.Namespace)
			for ri, r := range kvs.Reads {
				golog.Infof("  .KVRWSet.Reads[%d]: %#v", ri, r)
			}
			for wi, w := range kvs.Writes {
				golog.Infof("  .KVRWSet.Writes[%d]: %#v", wi, w)
			}
		}
	}

	return nil
}

// Handler export
var Handler handler
