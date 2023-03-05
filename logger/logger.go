package logger

import (
	"io"

	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000000000Z07:00"
}

func New(caller string, w io.Writer, lv zerolog.Level) zerolog.Logger {
	logger := zerolog.New(w).Level(lv).With().Timestamp().Logger()
	if caller != "" {
		logger = logger.With().Str("caller", caller).Logger()
	}
	return logger
}

func NewConsoleWriter(options ...func(w *zerolog.ConsoleWriter)) zerolog.ConsoleWriter {
	opts := append([]func(w *zerolog.ConsoleWriter){ // default
		func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = zerolog.TimeFieldFormat //time.RFC3339
		},
	}, options...)
	return zerolog.NewConsoleWriter(opts...)
}
