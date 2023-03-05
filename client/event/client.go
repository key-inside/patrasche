package event

import (
	"errors"
	"fmt"
	"time"

	fabopts "github.com/hyperledger/fabric-sdk-go/pkg/common/options"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/client"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/service/dispatcher"
)

type Client struct {
	eventService fab.EventService
}

type params struct {
	seekType                seek.Type
	fromBlock               uint64
	eventConsumerBufferSize uint
	eventConsumerTimeout    time.Duration
}

// New returns a client instance permits block events (WithBlockEvents) and has custom backpressure strategy (WithEventConsumerTimeout(0))
func New(channelProvider context.ChannelProvider, options ...Option) (*Client, error) {
	channelContext, err := channelProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel context: %w", err)
	}

	if channelContext.ChannelService() == nil {
		return nil, errors.New("channel service not initialized")
	}

	// default option parameters
	// IMPORTANT: if eventConsumerTimeout not 0, blocks can be omitted when the event channel buffer is full.
	//            or you can set eventConsumerBufferSize to 0
	p := &params{
		seekType:                seek.Newest,
		fromBlock:               0,
		eventConsumerBufferSize: 100,
		eventConsumerTimeout:    0,
	}
	for _, option := range options {
		option(p)
	}

	opts := []fabopts.Opt{
		client.WithBlockEvents(),
		deliverclient.WithSeekType(p.seekType),
		deliverclient.WithBlockNum(p.fromBlock),
		dispatcher.WithEventConsumerBufferSize(p.eventConsumerBufferSize),
		dispatcher.WithEventConsumerTimeout(p.eventConsumerTimeout),
	}

	es, err := channelContext.ChannelService().EventService(opts...)
	if err != nil {
		return nil, fmt.Errorf("event service creation failed: %w", err)
	}

	return &Client{eventService: es}, nil
}

// RegisterBlockEvent registers for block events. If the caller does not have permission
// to register for block events then an error is returned. Unregister must be called when the registration is no longer needed.
func (c *Client) RegisterBlockEvent(filter ...fab.BlockFilter) (fab.Registration, <-chan *fab.BlockEvent, error) {
	return c.eventService.RegisterBlockEvent(filter...)
}

// RegisterFilteredBlockEvent registers for filtered block events. Unregister must be called when the registration is no longer needed.
func (c *Client) RegisterFilteredBlockEvent() (fab.Registration, <-chan *fab.FilteredBlockEvent, error) {
	return c.eventService.RegisterFilteredBlockEvent()
}

// RegisterChaincodeEvent registers for chaincode events. Unregister must be called when the registration is no longer needed.
func (c *Client) RegisterChaincodeEvent(ccID, eventFilter string) (fab.Registration, <-chan *fab.CCEvent, error) {
	return c.eventService.RegisterChaincodeEvent(ccID, eventFilter)
}

// RegisterTxStatusEvent registers for transaction status events. Unregister must be called when the registration is no longer needed.
func (c *Client) RegisterTxStatusEvent(txID string) (fab.Registration, <-chan *fab.TxStatusEvent, error) {
	return c.eventService.RegisterTxStatusEvent(txID)
}

// Unregister removes the given registration and closes the event channel.
func (c *Client) Unregister(reg fab.Registration) {
	c.eventService.Unregister(reg)
}

type Option func(p *params)

func WithBlockNum(value uint64) Option {
	return func(p *params) {
		p.seekType = seek.FromBlock
		p.fromBlock = value
	}
}

// WithEventConsumerBufferSize sets the size of the registered consumer's event channel.
func WithEventConsumerBufferSize(value uint) Option {
	return func(p *params) {
		p.eventConsumerBufferSize = value
	}
}

// WithEventConsumerTimeout is the timeout when sending events to a registered consumer.
// If < 0, if buffer full, unblocks immediately and does not send.
// If 0, if buffer full, will block and guarantee the event will be sent out.
// If > 0, if buffer full, blocks util timeout.
func WithEventConsumerTimeout(value time.Duration) Option {
	return func(p *params) {
		p.eventConsumerTimeout = value
	}
}
