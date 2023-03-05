package tx

import (
	"regexp"

	"github.com/hyperledger/fabric-protos-go/common"
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

func (f *hashFilter) Handle(tx *Tx) error {
	if f.hashPattern.MatchString(tx.ID()) {
		return f.next.Handle(tx)
	}
	for _, action := range f.filteredActions {
		if err := action(tx); err != nil {
			return err
		}
	}
	return nil
}

func NewHashFilteredLoggingAction(logger *zerolog.Logger) Action {
	return func(tx *Tx) error {
		logger.Debug().
			Uint64("block_number", tx.BlockNum).
			Str("id", tx.ID()).
			Msg("tx filtered by hash")
		return nil
	}
}

type validEndorserFilter struct {
	next            Handler
	filteredActions []Action
}

func NewValidEndorserFilter(next Handler, filteredActions ...Action) Handler {
	return &validEndorserFilter{
		next:            next,
		filteredActions: filteredActions,
	}
}

func (f *validEndorserFilter) Handle(tx *Tx) error {
	if tx.IsValid() && tx.HeaderType() == common.HeaderType_ENDORSER_TRANSACTION {
		return f.next.Handle(tx)
	}
	for _, action := range f.filteredActions {
		if err := action(tx); err != nil {
			return err
		}
	}
	return nil
}

func NewValidEndorserFilteredLoggingAction(logger *zerolog.Logger) Action {
	return func(tx *Tx) error {
		logger.Debug().
			Uint64("block_number", tx.BlockNum).
			Str("id", tx.ID()).
			Str("type", tx.HeaderType().String()).
			Str("validation", tx.ValidationCode.String()).
			Msg("tx filtered by valid-endorser")
		return nil
	}
}
