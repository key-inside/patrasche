package ledger

import (
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/key-inside/patrasche"
	"github.com/key-inside/patrasche/block"
	"github.com/key-inside/patrasche/tx"
)

var once sync.Once

var cmd *cobra.Command

func Command() *cobra.Command {
	once.Do(func() {
		cmd = &cobra.Command{
			Use:   "ldg",
			Short: "Query Ledger Data",
			Long:  "Querying blocks and transactions from the ledger",
			Run: func(cmd *cobra.Command, args []string) {
				p := patrasche.Biter(cmd)
				logger := p.Logger().With().Str("caller", "ldg").Logger()

				ch, err := p.NewChannel()
				if err != nil {
					logger.Error().Err(err).Send()
					return
				}
				defer ch.Close()

				client, err := ch.NewLedgerClient()
				if err != nil {
					logger.Error().Err(err).Send()
					return
				}

				if viper.IsSet("block") { // query block
					bn := viper.GetUint64("block")
					b, err := client.QueryBlock(bn)
					if err != nil {
						logger.Error().Err(err).Send()
						return
					}
					blockLogger := block.NewStdLogger(nil, &logger)
					if err := blockLogger.Handle(b); err != nil {
						logger.Error().Err(err).Send()
						return
					}
				} else if viper.IsSet("txid") { // query tx
					txID := strings.ToLower(viper.GetString("txid"))
					t, err := client.QueryTransaction(txID)
					if err != nil {
						logger.Error().Err(err).Send()
						return
					}
					txLogger := tx.NewStdLogger(nil, &logger)
					if err := txLogger.Handle(t); err != nil {
						logger.Error().Err(err).Send()
						return
					}
				}
			},
		}

		flags := cmd.Flags()
		flags.Uint64P("block", "b", 0, "block number")
		flags.StringP("txid", "t", "", "tx ID (hex)")

		viper.BindPFlags(flags)
	})

	return cmd
}
