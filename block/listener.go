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
	defer client.Unregister(registration)

	var wg sync.WaitGroup

	errCh := make(chan error, 1)
	stopCh := make(chan bool, 1)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		stopCh <- true
		wg.Wait()
		if l.shutdown != nil {
			l.shutdown(sig)
		}
		os.Exit(1)
	}()

	wg.Add(1)
	go func() {
	LISTENING:
		for {
			select {
			case evt := <-notifier: // wait block event
				block, err := New(evt.Block)
				if err != nil {
					errCh <- fmt.Errorf("failed to parse block data: %w", err)
					break LISTENING
				}
				if err := l.handler.Handle(block); err != nil {
					errCh <- err
					break LISTENING
				}
				if l.endBlock != nil && block.Num >= *l.endBlock {
					errCh <- nil
					break LISTENING
				}
			case <-stopCh:
				break LISTENING
			}
		}
		wg.Done()
	}()

	return <-errCh
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
