// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package cmd

import (
	"github.com/kataras/golog"
	"github.com/spf13/cobra"

	"github.com/key-inside/patrasche/handler/inspect"
	"github.com/key-inside/patrasche/pkg/listener"
)

func init() {
	// add it to root command
	rootCmd.AddCommand(inspectCmd)
}

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect transactions",
	Long:  "Inspect transactions (ENDORSER_TRANSACTION)",
	Run: func(cmd *cobra.Command, args []string) {
		if err := listener.Listen(inspect.Handler); err != nil {
			golog.Fatal(err)
		}
	},
}
