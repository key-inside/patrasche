package block

import (
	"os"
	"strconv"

	"github.com/key-inside/patrasche/aws"
)

type Handler interface {
	Handle(block *Block) error
}

type blockNumerFileWriter struct {
	path string
	next Handler
}

func NewBlockNumberFileWriter(next Handler, path string) Handler {
	return &blockNumerFileWriter{
		path: path,
		next: next,
	}
}

func (h *blockNumerFileWriter) Handle(block *Block) error {
	err := os.WriteFile(h.path, []byte(strconv.FormatUint(block.Num, 10)), 0644)
	if err != nil {
		return err
	}
	if h.next != nil {
		return h.next.Handle(block)
	}
	return nil
}

type blockNumerDynamoDBWriter struct {
	table       string
	itemFactory func(*Block) any
	next        Handler
}

func NewBlockNumberDynamoDBWriter(next Handler, table string, itemFactory func(*Block) any) Handler {
	return &blockNumerDynamoDBWriter{
		table:       table,
		itemFactory: itemFactory,
		next:        next,
	}
}

func (h *blockNumerDynamoDBWriter) Handle(block *Block) error {
	return aws.PutItemToDynamoDB(aws.DefaultConfig(), h.table, h.itemFactory(block))
}
