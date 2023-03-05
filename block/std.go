package block

import (
	"errors"

	"github.com/rs/zerolog"

	"github.com/key-inside/patrasche/tx"
)

type stdHandler struct {
	handler tx.Handler
}

func NewStdHandler(handler tx.Handler) (Handler, error) {
	if handler == nil {
		return nil, errors.New("tx handler is nil")
	}

	return &stdHandler{handler: handler}, nil
}

func (h *stdHandler) Handle(block *Block) error {
	if block != nil {
		for _, t := range block.Txs {
			if err := h.handler.Handle(t); err != nil {
				return err
			}
		}
	}
	return nil
}

type stdLogger struct {
	logger *zerolog.Logger
	next   Handler
}

func NewStdLogger(next Handler, logger *zerolog.Logger) Handler {
	return &stdLogger{
		logger: logger,
		next:   next,
	}
}

func (h *stdLogger) Handle(block *Block) error {
	h.logger.Info().
		Uint64("number", block.Num).
		Hex("hash", block.Hash).
		Int("tx_count", len(block.Txs)).
		Msg("block")
	if h.next != nil {
		return h.next.Handle(block)
	}
	return nil
}
