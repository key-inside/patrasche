package inspect

import (
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/rs/zerolog"

	"github.com/key-inside/patrasche/tx"
)

type inspectHandler struct {
	logger zerolog.Logger
}

func NewTxHandler(logger zerolog.Logger) tx.Handler {
	return &inspectHandler{logger: logger}
}

func (h *inspectHandler) Handle(t *tx.Tx) error {
	dic := zerolog.Dict().
		Uint64("block_num", t.BlockNum).
		Str("id", t.ID()).
		Str("mspid", t.MSPID()).
		Str("validation_code", t.ValidationCode.String()).
		Str("type", t.HeaderType().String())

	if t.HeaderType() == common.HeaderType_ENDORSER_TRANSACTION {
		// ChaincodeInvocationSpec
		spec, err := t.GetChaincodeInvocationSpec()
		if err != nil {
			return err
		}
		subDic := zerolog.Dict().Any("chaincode_id", spec.ChaincodeSpec.ChaincodeId)
		args := []string{}
		for _, arg := range spec.ChaincodeSpec.Input.Args {
			args = append(args, string(arg))
		}
		subDic.Any("args", args)
		dic.Dict("chaincode_invocation_spec", subDic)

		// ChaincodeAction
		ccA, err := t.GetChaincodeAction()
		if err != nil {
			return err
		}
		subDic = zerolog.Dict().
			Any("chaincode_id", ccA.ChaincodeId).
			Dict("response", zerolog.Dict().
				Int("status", int(ccA.Response.Status)).
				Str("payload", string(ccA.Response.Payload)).
				Str("message", ccA.Response.Message))

		// ChaincodeEvent
		ccE, err := t.GetChaincodeEvent()
		if err != nil {
			return err
		}
		if ccE != nil {
			subDic.Dict("event", zerolog.Dict().
				Str("name", ccE.EventName).
				Str("payload", string(ccE.Payload)))
		}
		dic.Any("chaincode_action", subDic)

		// RWSet
		rwm, err := t.GetReadWriteMap()
		if err != nil {
			return err
		}
		dic.Any("rwset", rwm)
	}

	h.logger.Info().Dict("tx", dic).Msg("")

	return nil
}
