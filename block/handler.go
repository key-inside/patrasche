package block

import (
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"

	patrasche_aws "github.com/key-inside/patrasche/aws"
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
	cfg         aws.Config
	table       string
	itemFactory func(*Block) any
	next        Handler
}

func NewBlockNumberDynamoDBWriter(next Handler, awsCfg aws.Config, table string, itemFactory func(*Block) any) Handler {
	return &blockNumerDynamoDBWriter{
		cfg:         awsCfg,
		table:       table,
		itemFactory: itemFactory,
		next:        next,
	}
}

func (h *blockNumerDynamoDBWriter) Handle(block *Block) error {
	err := patrasche_aws.PutItemToDynamoDB(h.cfg, h.table, h.itemFactory(block))
	if err != nil {
		return err
	}
	if h.next != nil {
		return h.next.Handle(block)
	}
	return nil
}
