package block

import (
	"encoding/hex"
	"regexp"

	"github.com/rs/zerolog"
)

type hashFilter struct {
	hashPattern     *regexp.Regexp
	next            Handler
	filteredActions []Action
}

func NewHashFilter(next Handler, pattern string, filteredActions ...Action) Handler {
	return &hashFilter{
		hashPattern:     regexp.MustCompilePOSIX(pattern),
		next:            next,
		filteredActions: filteredActions,
	}
}

func (f *hashFilter) Handle(block *Block) error {
	if f.hashPattern.MatchString(hex.EncodeToString(block.Hash)) {
		return f.next.Handle(block)
	}
	for _, action := range f.filteredActions {
		if err := action(block); err != nil {
			return err
		}
	}
	return nil
}

func NewHashFilteredLoggingAction(logger *zerolog.Logger) Action {
	return func(block *Block) error {
		logger.Debug().
			Uint64("number", block.Num).
			Hex("hash", block.Hash).
			Msg("block filtered by hash")
		return nil
	}
}
