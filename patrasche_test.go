// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package patrasche

import (
	"testing"

	"github.com/key-inside/patrasche/cmd/inspect"
)

func Test_Inspect(t *testing.T) {
	cmd := NewRootCommand(&App{
		Name: "test",
	})
	cmd.AddCommand(inspect.Command())
	cmd.Execute()
}
