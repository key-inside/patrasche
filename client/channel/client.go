package channel

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
)

type Client struct {
	*channel.Client
}

func New(ctx context.ChannelProvider, options ...channel.ClientOption) (*Client, error) {
	_client, err := channel.New(ctx, options...)
	if err != nil {
		return nil, err
	}
	return &Client{_client}, nil
}
