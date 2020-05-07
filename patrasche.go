// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package patrasche

import (
	"fmt"

	"github.com/kataras/golog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/key-inside/patrasche/pkg/aws"
	"github.com/key-inside/patrasche/pkg/listener"
	"github.com/key-inside/patrasche/pkg/logo"
	"github.com/key-inside/patrasche/pkg/tx"
	"github.com/key-inside/patrasche/pkg/version"
)

const fixedRFC3339Nano = "2006-01-02T15:04:05.000000000Z07:00"

func init() {
	golog.SetTimeFormat(fixedRFC3339Nano)
}

// App options for creating the root command
type App struct {
	Name      string
	Short     string
	Long      string
	EnvPrefix string
	Version   string
}

// NewRootCommand _
func NewRootCommand(app *App) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     app.Name,
		Short:   app.Short,
		Long:    app.Long,
		Version: version.Version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logo.Print()
			fmt.Println(cmd.VersionTemplate())
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			fmt.Println()
		},
	}

	rootCmd.SetVersionTemplate(version.TemplatedVersion)

	// patrasche config
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file path or ARN")
	rootCmd.PersistentFlags().String("config.region", "", "config SSM region")
	rootCmd.PersistentFlags().String("config.parameter", "", "config SSM parameter")
	// network config
	rootCmd.PersistentFlags().StringP("network", "n", "", "network config file path or ARN")
	rootCmd.PersistentFlags().String("network.region", "", "network SSM region")
	rootCmd.PersistentFlags().String("network.parameter", "", "network SSM parameter")
	// network identity
	rootCmd.PersistentFlags().StringP("identity", "u", "", "user id, msp path or ARN")
	rootCmd.PersistentFlags().String("identity.region", "", "identity SSM region")
	rootCmd.PersistentFlags().String("identity.parameter", "", "identity SSM parameter")
	// network channel
	rootCmd.PersistentFlags().StringP("channel", "C", "", "channel name")
	// block
	rootCmd.PersistentFlags().StringP("block", "b", "", "number or file path or ARN for block number, unset or empty is newest block")
	rootCmd.PersistentFlags().String("block.region", "", "block SSM region")
	rootCmd.PersistentFlags().String("block.parameter", "", "block SSM parameter")
	// follow next blocks
	rootCmd.PersistentFlags().BoolP("follow", "f", false, "follow event")
	// tx filter
	rootCmd.PersistentFlags().StringP("tx.id", "i", "", "tx ID pattern")
	rootCmd.PersistentFlags().StringP("tx.type", "t", "", "tx header type")
	// logging
	rootCmd.PersistentFlags().String("logging.level", "info", "fatal, error, warn, info, debug or disable")

	viper.BindPFlags(rootCmd.PersistentFlags())
	if app.EnvPrefix != "" {
		viper.SetEnvPrefix(app.EnvPrefix)
	} else {
		viper.SetEnvPrefix("PATRASCHE")
	}
	viper.AutomaticEnv()

	cobra.OnInitialize(initConfig)

	return rootCmd
}

func initConfig() {
	// set config to viper
	arn, path, err := aws.GetARN("config")
	if err != nil {
		if path != "" {
			viper.SetConfigFile(path)
			if err := viper.ReadInConfig(); err != nil {
				golog.Fatal("Can't read config:", err)
			}
		}
	} else { // AWS resource
		in, typ, err := aws.GetReaderWithARN(arn)
		if err != nil {
			golog.Fatal("Can't get config from AWS:", err)
		}
		viper.SetConfigType(typ)
		if err := viper.ReadConfig(in); err != nil {
			golog.Fatal("Can't read config:", err)
		}
	}

	// update logging level
	if lv := viper.GetString("logging.level"); lv != "" {
		golog.SetLevel(lv)
	}
}

// Listen listens block events
func Listen(txh tx.Handler) error {
	return listener.Listen(txh)
}
