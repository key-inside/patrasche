package block

import (
	"github.com/hyperledger/fabric-protos-go/common"

	"github.com/key-inside/patrasche/proto"
	"github.com/key-inside/patrasche/tx"
)

type Block struct {
	*common.Block

	Num  uint64
	Hash []byte
	Txs  []*tx.Tx
}

func New(block *common.Block) (*Block, error) {
	hash, err := GenerateHash(block)
	if err != nil {
		return nil, err
	}

	txs := []*tx.Tx{}
	for i, data := range block.Data.Data {
		validationByte := block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER][i] // BlockMetadataIndex_TRANSACTIONS_FILTER = 2
		envelope, err := proto.UnmarshalEnvelope(data)
		if err != nil {
			return nil, err
		}
		t, err := tx.New(block.Header.Number, i, validationByte, envelope.Payload)
		if err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}

	return &Block{
		Block: block,
		Num:   block.Header.Number,
		Hash:  hash,
		Txs:   txs,
	}, nil
}
