package inspect

import (
	"os"
	"strconv"
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
			Use:     "inspect",
			Short:   "Inspect transactions",
			Long:    "Inspect transactions",
			Version: "v0.0.0",
			Run: func(cmd *cobra.Command, args []string) {
				p := patrasche.Biter(cmd)
				logger := p.Logger().With().Str("caller", "inspect").Logger()

				// inspect tx handler
				txHandler := NewTxHandler(logger)
				// tx filters
				if pattern := viper.GetString("filter.tx-hash"); pattern != "" {
					txHandler = tx.NewHashFilter(txHandler, pattern, tx.NewHashFilteredLoggingAction(&logger))
				}
				if viper.GetBool("filter.valid-endorser") {
					txHandler = tx.NewValidEndorserFilter(txHandler, tx.NewValidEndorserFilteredLoggingAction(&logger))
				}
				// logging middleware
				txHandler = tx.NewStdLogger(txHandler, &logger)

				// standard block handler
				blockHandler, err := block.NewStdHandler(txHandler)
				if err != nil {
					logger.Error().Err(err).Msg("")
					return
				}
				// block filter
				if pattern := viper.GetString("filter.block-hash"); pattern != "" {
					blockHandler = block.NewHashFilter(blockHandler, pattern, block.NewHashFilteredLoggingAction(&logger))
				}
				// logging middleware
				blockHandler = block.NewStdLogger(blockHandler, &logger)

				// listener options
				opts := []block.ListenerOption{}
				var startBn uint64
				if path := viper.GetString("save"); path != "" {
					nb, _ := os.ReadFile(path)
					startBn, _ = strconv.ParseUint(string(nb), 10, 64)
					opts = append(opts, block.WithStartBlock(startBn))
					blockHandler = block.NewBlockNumberFileWriter(blockHandler, path) // save block number before block filtering
				}
				if viper.IsSet("start") {
					bn := viper.GetUint64("start")
					if startBn < bn {
						opts = append(opts, block.WithStartBlock(bn)) // it will override previous WithStartBlock value
					}
				}
				if viper.IsSet("end") {
					opts = append(opts, block.WithEndBlock(viper.GetUint64("end")))
				}

				if err := p.ListenBlock(blockHandler, opts...); err != nil {
					logger.Error().Err(err).Msg("")
					return
				}

				/*
					// alternative style
					ch, err := p.NewChannel(fabsdk.WithUser("nello"))
					if err != nil {
						logger.Error().Err(err).Msg("")
						return
					}
					defer ch.Close()
					l, err := block.NewListener(ch, blockHandler, opts...)
					if err != nil {
						logger.Error().Err(err).Msg("")
						return
					}
					if err := l.Listen(); err != nil {
						logger.Error().Err(err).Msg("")
						return
					}
				*/
			},
		}

		flags := cmd.Flags()
		flags.String("save", "", "block number save file path")
		flags.Uint64("start", 0, "start block number, if not set, seek from newest")
		flags.Uint64("end", 0, "end block number")
		flags.Bool("filter.valid-endorser", false, "valid endorser tx only")
		flags.String("filter.block-hash", "", "block hash pattern")
		flags.String("filter.tx-hash", "", "tx hash pattern")

		viper.BindPFlags(flags)
	})

	return cmd
}
