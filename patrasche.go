// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package patrasche

import (
	"github.com/spf13/cobra"

	"github.com/key-inside/patrasche/cmd"
	"github.com/key-inside/patrasche/pkg/listener"
	"github.com/key-inside/patrasche/pkg/tx"
)

// AddCommand adds a sub command to the root command
func AddCommand(sub *cobra.Command) {
	cmd.AddCommand(sub)
}

// Execute executes the root command
func Execute() error {
	return cmd.Execute()
}

// Listen listens block events
func Listen(txh tx.Handler) error {
	return listener.Listen(txh)
}
