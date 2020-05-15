// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package patrasche

import (
	"fmt"
	"strings"

	"github.com/kataras/golog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/key-inside/patrasche/pkg/aws"
	"github.com/key-inside/patrasche/pkg/flag"
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
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			fmt.Println()
		},
	}

	rootCmd.SetVersionTemplate(fmt.Sprintf("%s %s (%s)\n", app.Name, app.Version, version.TemplatedVersion))

	pFlags := rootCmd.PersistentFlags()

	// patrasche config
	pFlags.AddFlagSet(flag.NewARNFlagSet("config", "", "config source (path, ARN)"))
	// network config
	pFlags.AddFlagSet(flag.NewARNFlagSet("patrasche.network", "", "network config source (path, ARN)"))
	// network identity
	pFlags.AddFlagSet(flag.NewARNFlagSet("patrasche.identity", "", "identity (user ID, msp path, ARN)"))
	// network channel
	pFlags.String("patrasche.channel", "", "channel name")
	// block
	pFlags.AddFlagSet(flag.NewARNFlagSet("patrasche.block", "", "block number (number, path, ARN), unset or empty is newest block"))
	// follow next blocks
	pFlags.Bool("patrasche.follow", false, "follow event")
	// tx filter
	pFlags.String("patrasche.tx.id", "", "tx ID pattern")
	pFlags.String("patrasche.tx.type", "", "tx header type")
	// logging
	pFlags.String("patrasche.logging.level", "info", "fatal, error, warn, info, debug or disable")

	viper.BindPFlags(pFlags)
	if app.EnvPrefix != "" {
		viper.SetEnvPrefix(app.EnvPrefix)
	} else {
		viper.SetEnvPrefix(app.Name)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

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
				golog.Fatal("Can't read config: ", err)
			}
		}
	} else { // AWS resource
		in, typ, err := aws.GetReaderWithARN(arn)
		if err != nil {
			golog.Fatal("Can't get config from AWS: ", err)
		}
		viper.SetConfigType(typ)
		if err := viper.ReadConfig(in); err != nil {
			golog.Fatal("Can't read config: ", err)
		}
	}

	// update logging level
	if lv := viper.GetString("patrasche.logging.level"); lv != "" {
		golog.SetLevel(lv)
	}
}

// Listen listens block events
func Listen(txh tx.Handler) error {
	return listener.Listen(txh)
}
