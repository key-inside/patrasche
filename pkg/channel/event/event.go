// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package event

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/common/options"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/client"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/service/dispatcher"
	"github.com/pkg/errors"
)

// Client _
type Client struct {
	eventService fab.EventService
}

// New _
func New(channelProvider context.ChannelProvider, fromBlock uint64, seekType seek.Type) (*Client, error) {
	channelContext, err := channelProvider()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create channel context")
	}

	if channelContext.ChannelService() == nil {
		return nil, errors.New("channel service not initialized")
	}

	// IMPORTANT: if eventConsumerTimeout not 0, blocks can be omitted when the event channel buffer is full.
	//            or you can set eventConsumerBufferSize to 0
	opts := []options.Opt{
		dispatcher.WithEventConsumerTimeout(0),
		// dispatcher.WithEventConsumerBufferSize(0),
		client.WithBlockEvents(),
	}
	if seekType != "" {
		opts = append(opts, deliverclient.WithSeekType(seekType))
		if seekType == seek.FromBlock {
			opts = append(opts, deliverclient.WithBlockNum(fromBlock))
		}
	}

	es, err := channelContext.ChannelService().EventService(opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "event service creation failed")
	}

	return &Client{eventService: es}, nil
}

// RegisterBlockEvent registers for block events. If the caller does not have permission
// to register for block events then an error is returned. Unregister must be called when the registration is no longer needed.
func (c *Client) RegisterBlockEvent(filter ...fab.BlockFilter) (fab.Registration, <-chan *fab.BlockEvent, error) {
	return c.eventService.RegisterBlockEvent(filter...)
}

// // RegisterFilteredBlockEvent registers for filtered block events. Unregister must be called when the registration is no longer needed.
// func (c *Client) RegisterFilteredBlockEvent() (fab.Registration, <-chan *fab.FilteredBlockEvent, error) {
// 	return c.eventService.RegisterFilteredBlockEvent()
// }

// // RegisterChaincodeEvent registers for chaincode events. Unregister must be called when the registration is no longer needed.
// func (c *Client) RegisterChaincodeEvent(ccID, eventFilter string) (fab.Registration, <-chan *fab.CCEvent, error) {
// 	return c.eventService.RegisterChaincodeEvent(ccID, eventFilter)
// }

// // RegisterTxStatusEvent registers for transaction status events. Unregister must be called when the registration is no longer needed.
// func (c *Client) RegisterTxStatusEvent(txID string) (fab.Registration, <-chan *fab.TxStatusEvent, error) {
// 	return c.eventService.RegisterTxStatusEvent(txID)
// }

// Unregister removes the given registration and closes the event channel.
func (c *Client) Unregister(reg fab.Registration) {
	c.eventService.Unregister(reg)
}
