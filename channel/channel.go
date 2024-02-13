package channel

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	chclient "github.com/key-inside/patrasche/client/channel"
	evtclient "github.com/key-inside/patrasche/client/event"
	ldgclient "github.com/key-inside/patrasche/client/ledger"
)

type Channel struct {
	sdk   *fabsdk.FabricSDK
	chCtx context.ChannelProvider
}

func New(channelID string, ctx core.ConfigProvider, ctxOpts ...fabsdk.ContextOption) (*Channel, error) {
	sdk, err := fabsdk.New(ctx) // sdk, err := fabsdk.New(ctx, ([]fabsdk.Option{})...)
	if err != nil {
		return nil, err
	}

	chCtx := sdk.ChannelContext(channelID, ctxOpts...)

	return &Channel{
		sdk:   sdk,
		chCtx: chCtx,
	}, nil
}

func (c *Channel) Close() {
	if c.sdk != nil {
		c.sdk.Close()
	}
}

func (c *Channel) NewClient(options ...channel.ClientOption) (*chclient.Client, error) {
	return chclient.New(c.chCtx, options...)
}

func (c *Channel) NewBlockEventClient(options ...evtclient.Option) (*evtclient.Client, error) {
	return evtclient.New(c.chCtx, options...)
}

func (c *Channel) NewLedgerClient(options ...ledger.ClientOption) (*ldgclient.Client, error) {
	return ldgclient.New(c.chCtx, options...)
}
