package ledger

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"

	"github.com/key-inside/patrasche/block"
	"github.com/key-inside/patrasche/tx"
)

type Client struct {
	*ledger.Client
}

func New(ctx context.ChannelProvider, options ...ledger.ClientOption) (*Client, error) {
	_client, err := ledger.New(ctx, options...)
	if err != nil {
		return nil, err
	}
	return &Client{_client}, nil
}

func (c *Client) QueryBlock(blockNumber uint64, options ...ledger.RequestOption) (*block.Block, error) {
	b, err := c.Client.QueryBlock(blockNumber, options...)
	if err != nil {
		return nil, err
	}
	return block.New(b)
}

func (c *Client) QueryBlockByHash(blockHash []byte, options ...ledger.RequestOption) (*block.Block, error) {
	b, err := c.Client.QueryBlockByHash(blockHash, options...)
	if err != nil {
		return nil, err
	}
	return block.New(b)
}

func (c *Client) QueryBlockByTxID(txID string, options ...ledger.RequestOption) (*block.Block, error) {
	b, err := c.Client.QueryBlockByTxID(fab.TransactionID(txID), options...)
	if err != nil {
		return nil, err
	}
	return block.New(b)
}

func (c *Client) QueryTransaction(txID string, options ...ledger.RequestOption) (*tx.Tx, error) {
	b, err := c.QueryBlockByTxID(txID, options...)
	if err != nil {
		return nil, err
	}
	for _, t := range b.Txs {
		if t.ID() == txID {
			return t, nil
		}
	}
	return nil, fmt.Errorf("tx disappeared") // never here
}

func (c *Client) FastQueryTransaction(txID string, options ...ledger.RequestOption) (*tx.Tx, error) {
	pbTx, err := c.Client.QueryTransaction(fab.TransactionID(txID), options...)
	if err != nil {
		return nil, err
	}
	return tx.New(0, 0, byte(pbTx.ValidationCode), pbTx.TransactionEnvelope.Payload)
}

type BlockchainInfoResponse struct {
	*fab.BlockchainInfoResponse
}

func (c *Client) QueryInfo(options ...ledger.RequestOption) (*BlockchainInfoResponse, error) {
	bir, err := c.Client.QueryInfo(options...)
	if err != nil {
		return nil, err
	}
	return &BlockchainInfoResponse{bir}, nil
}

type ChannelCfg struct {
	fab.ChannelCfg
}

func (c *Client) QueryConfig(options ...ledger.RequestOption) (*ChannelCfg, error) {
	cfg, err := c.Client.QueryConfig(options...)
	if err != nil {
		return nil, err
	}
	return &ChannelCfg{cfg}, nil
}

func (c *Client) QueryConfigBlock(options ...ledger.RequestOption) (*block.Block, error) {
	b, err := c.Client.QueryConfigBlock(options...)
	if err != nil {
		return nil, err
	}
	return block.New(b)
}
