package ccquery

import (
	"sync"

	fabch "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/key-inside/patrasche"
)

var once sync.Once

var cmd *cobra.Command

func Command() *cobra.Command {
	once.Do(func() {
		cmd = &cobra.Command{
			Use:   "ccq",
			Short: "Chaincode query",
			Long:  "Querying chaincode function retrieves data",
			Run: func(cmd *cobra.Command, args []string) {
				p := patrasche.Biter(cmd)
				logger := p.Logger().With().Str("caller", "ccq").Logger()

				ch, err := p.NewChannel()
				if err != nil {
					logger.Error().Err(err).Send()
					return
				}
				defer ch.Close()

				client, err := ch.NewClient()
				if err != nil {
					logger.Error().Err(err).Send()
					return
				}

				reqArgs := [][]byte{}
				for _, arg := range viper.GetStringSlice("args") {
					reqArgs = append(reqArgs, []byte(arg))
				}

				req := fabch.Request{
					ChaincodeID: viper.GetString("cc"),
					Fcn:         viper.GetString("fn"),
					Args:        reqArgs,
				}
				res, err := client.Query(req)
				if err != nil {
					logger.Error().Err(err).Send()
					return
				}
				logger.Info().Str("payload", string(res.Payload)).Msg("query success")
			},
		}

		flags := cmd.Flags()
		flags.String("cc", "", "chaincode name")
		flags.String("fn", "", "function name")
		flags.StringArray("args", []string{}, "arguments")

		viper.BindPFlags(flags)
	})

	return cmd
}
