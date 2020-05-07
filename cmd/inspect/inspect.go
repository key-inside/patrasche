// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package inspect

import (
	"github.com/kataras/golog"
	"github.com/spf13/cobra"

	"github.com/key-inside/patrasche/handler/inspect"
	"github.com/key-inside/patrasche/pkg/listener"
)

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

// Command _
func Command() *cobra.Command {
	return inspectCmd
}
