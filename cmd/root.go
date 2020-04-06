// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package cmd

import (
	"fmt"

	"github.com/kataras/golog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/key-inside/patrasche/pkg/config"
	"github.com/key-inside/patrasche/pkg/logo"
	"github.com/key-inside/patrasche/pkg/version"
)

const fixedRFC3339Nano = "2006-01-02T15:04:05.000000000Z07:00"

func init() {
	golog.SetTimeFormat(fixedRFC3339Nano)

	cobra.OnInitialize(initConfig)

	rootCmd.SetVersionTemplate(version.TemplatedVersion)

	rootCmd.PersistentFlags().String("config", "", "config file path or ARN")
	rootCmd.PersistentFlags().String("config.region", "", "config SSM region")
	rootCmd.PersistentFlags().String("config.parameter", "", "config SSM parameter")
	rootCmd.PersistentFlags().StringP("network", "n", "", "network config file path or ARN")
	rootCmd.PersistentFlags().String("network.region", "", "network SSM region")
	rootCmd.PersistentFlags().String("network.parameter", "", "network SSM parameter")
	rootCmd.PersistentFlags().StringP("identity", "u", "", "user id, msp path or ARN")
	rootCmd.PersistentFlags().String("identity.region", "", "identity SSM region")
	rootCmd.PersistentFlags().String("identity.parameter", "", "identity SSM parameter")
	rootCmd.PersistentFlags().StringP("channel", "c", "", "channel name")
	rootCmd.PersistentFlags().StringP("block", "b", "", "number or file path or ARN for block number, unset or empty is newest block")
	rootCmd.PersistentFlags().String("block.region", "", "block SSM region")
	rootCmd.PersistentFlags().String("block.parameter", "", "block SSM parameter")
	rootCmd.PersistentFlags().BoolP("follow", "f", false, "follow event")
	rootCmd.PersistentFlags().StringP("tx.id", "i", "", "tx ID pattern")
	rootCmd.PersistentFlags().StringP("tx.type", "t", "", "tx header type")
	rootCmd.PersistentFlags().String("logging.level", "info", "fatal, error, warn, info, debug or disable")

	viper.BindPFlags(rootCmd.PersistentFlags())
}

func initConfig() {
	viper.SetEnvPrefix("PATRASCHE")
	viper.AutomaticEnv()

	// set config to viper
	arn, path, err := config.GetARN("config")
	if err != nil {
		if path != "" {
			viper.SetConfigFile(path)
			if err := viper.ReadInConfig(); err != nil {
				golog.Fatal("Can't read config:", err)
			}
		}
	} else { // AWS resource
		in, typ, err := config.GetReaderWithARN(arn)
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

var rootCmd = &cobra.Command{
	Use:     "patrasche",
	Short:   "Hyperledger Fabric Event Listener",
	Long:    "Subscribe to Hyperledger Fabric block events and handle transactions",
	Version: version.Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logo.Print()
		version.Print()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		fmt.Println()
	},
}

// AddCommand _
func AddCommand(sub *cobra.Command) {
	rootCmd.AddCommand(sub)
}

// Execute _
func Execute() error {
	return rootCmd.Execute()
}
