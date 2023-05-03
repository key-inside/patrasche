package block

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/key-inside/patrasche/channel"
	evtclient "github.com/key-inside/patrasche/client/event"
)

type Listener struct {
	ch         *channel.Channel
	handler    Handler
	startBlock *uint64
	endBlock   *uint64
	shutdown   func(os.Signal)
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
	opts := []evtclient.Option{}
	if l.startBlock != nil {
		opts = append(opts, evtclient.WithBlockNum(*l.startBlock))
	}
	client, err := l.ch.NewBlockEventClient(opts...)
	if err != nil {
		return fmt.Errorf("failed to create event client: %w", err)
	}
	registration, notifier, err := client.RegisterBlockEvent()
	if err != nil {
		return fmt.Errorf("failed to register block event: %w", err)
	}

	quitCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	stopping := false

	var m sync.Mutex
	quit := func(retErr error) {
		m.Lock()
		defer m.Unlock()
		if stopping {
			return
		}
		stopping = true
		// asynchronous unregister
		go func() {
			client.Unregister(registration)
			quitCh <- retErr
		}()
	}

	for {
		select {
		case sig := <-sigCh:
			quit(nil)
			if l.shutdown != nil {
				l.shutdown(sig)
			}
		case evt := <-notifier: // block event
			if evt != nil {
				if !stopping {
					block, err := New(evt.Block)
					if err != nil {
						quit(fmt.Errorf("failed to parse block data: %w", err))
						break
					}
					if err := l.handler.Handle(block); err != nil {
						quit(err)
						break
					}
					if l.endBlock != nil && block.Num >= *l.endBlock {
						quit(nil)
					}
				}
				// else MUST consume events in buffer for closing the channel.
				// because, when unregister event channel, fabric-sdk-go write nil to channel.
				// so, if the channel buffer is full, it will be pended forever.
			} else {
				quit(nil)
			}
		case e := <-quitCh:
			return e
		}
	}
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

func WithShutdown(shutdown func(os.Signal)) ListenerOption {
	return func(l *Listener) error {
		l.shutdown = shutdown
		return nil
	}
}
