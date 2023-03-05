package block

import (
	"errors"
	"fmt"

	"github.com/key-inside/patrasche/channel"
	evtclient "github.com/key-inside/patrasche/client/event"
)

type Listener struct {
	ch         *channel.Channel
	handler    Handler
	startBlock *uint64
	endBlock   *uint64
}

type ListenerOption func(*Listener) error

func NewListener(ch *channel.Channel, handler Handler, options ...ListenerOption) (*Listener, error) {
	if ch == nil {
		return nil, errors.New("channel is nil")
	}
	if handler == nil {
		return nil, errors.New("block handler is nil")
	}

	l := &Listener{ch: ch, handler: handler}

	for _, option := range options {
		if err := option(l); err != nil {
			return nil, fmt.Errorf("failed to apply listener option: %w", err)
		}
	}

	return l, nil
}

func (l *Listener) Listen() error {
	var blockNumOption evtclient.Option = nil
	if l.startBlock != nil {
		blockNumOption = evtclient.WithBlockNum(*l.startBlock)
	}
	client, err := l.ch.NewBlockEventClient(blockNumOption)
	if err != nil {
		return fmt.Errorf("failed to create event client: %w", err)
	}
	registration, notifier, err := client.RegisterBlockEvent()
	if err != nil {
		return fmt.Errorf("failed to register block event: %w", err)
	}
	defer client.Unregister(registration)

	for {
		evt := <-notifier // wait block event

		block, err := New(evt.Block)
		if err != nil {
			return fmt.Errorf("failed to parse block data: %w", err)
		}

		if err := l.handler.Handle(block); err != nil {
			return err
		}

		if l.endBlock != nil && block.Num >= *l.endBlock {
			break
		}
	}

	return nil
}

func WithStartBlock(blockNum uint64) ListenerOption {
	return func(l *Listener) error {
		l.startBlock = &blockNum
		return nil
	}
}

func WithEndBlock(blockNum uint64) ListenerOption {
	return func(l *Listener) error {
		l.endBlock = &blockNum
		return nil
	}
}
