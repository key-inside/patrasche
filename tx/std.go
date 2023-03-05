package tx

import "github.com/rs/zerolog"

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

func (h *stdLogger) Handle(tx *Tx) error {
	h.logger.Info().
		Str("id", tx.ID()).
		Str("type", tx.HeaderType().String()).
		Str("validation", tx.ValidationCode.String()).
		Str("timestamp", tx.Timestamp().String()).
		Msg("tx")
	if h.next != nil {
		return h.next.Handle(tx)
	}
	return nil
}
