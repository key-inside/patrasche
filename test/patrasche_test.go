package test

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"

	"github.com/key-inside/patrasche"
	"github.com/key-inside/patrasche/cmd/ccquery"
	"github.com/key-inside/patrasche/cmd/inspect"
	"github.com/key-inside/patrasche/cmd/ledger"
)

const testConfigs = "--config=./fixtures/config.yaml"

func bufferdOutput(c *cobra.Command) (out *bytes.Buffer) {
	out = bytes.NewBuffer(nil)
	c.SetOut(out)
	c.SetErr(out)
	return
}

func testCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "test",
		Short:   "Patrasche Test",
		Long:    "Patrasche Test Command",
		Version: "v0.0.0",
	}
}

func Test_Help(t *testing.T) {
	c := testCommand()
	c.AddCommand(inspect.Command(), ccquery.Command(), ledger.Command())
	// sets output before biting
	out := bufferdOutput(c)

	_, err := patrasche.Bite(c)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	c.Execute()
	t.Log(out.String())
}

func Test_RootCmd(t *testing.T) {
	c := ccquery.Command()
	// sets output before biting
	out := bufferdOutput(c)

	_, err := patrasche.Bite(c)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	c.SetArgs([]string{"--help"})
	if err := c.Execute(); err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	t.Log(out.String())
}

func Test_Inspect(t *testing.T) {
	c := testCommand()
	c.AddCommand(inspect.Command())
	// sets output before biting
	out := bufferdOutput(c)

	_, err := patrasche.Bite(c, patrasche.WithConsoleLogWriter())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// c.SetArgs([]string{"inspect", "--start=0", testConfigs})
	c.SetArgs([]string{"inspect", "--start=24049", "--end=0", testConfigs})
	// c.SetArgs([]string{"inspect", "--start=24049", "--end=0", "--filter.block-hash=^[abc1234]", testConfigs})
	// c.SetArgs([]string{"inspect", "--start=24049", "--end=24061", "--filter.tx-hash=^[abcd1234]", "--filter.valid-endorser", testConfigs})
	if err := c.Execute(); err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	t.Log(out.String())
}

func Test_CCQuery(t *testing.T) {
	c := testCommand()
	c.AddCommand(ccquery.Command())
	// sets output before biting
	out := bufferdOutput(c)

	_, err := patrasche.Bite(c)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	c.SetArgs([]string{"ccq", "--cc=ping", "--fn=ver", testConfigs})
	if err := c.Execute(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	t.Log(out.String())
}

func Test_QueryBlock(t *testing.T) {
	c := testCommand()
	c.AddCommand(ledger.Command())
	// sets output before biting
	out := bufferdOutput(c)

	_, err := patrasche.Bite(c)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	c.SetArgs([]string{"ldg", "-b=24049", testConfigs})
	if err := c.Execute(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	t.Log(out.String())
}

func Test_QueryTransaction(t *testing.T) {
	c := testCommand()
	c.AddCommand(ledger.Command())
	// sets output before biting
	out := bufferdOutput(c)

	_, err := patrasche.Bite(c)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	c.SetArgs([]string{"ldg", "-t=5e9ba4a754493424959bfeb52496652b99fa5be098d27a8e23c887fcff45de24", testConfigs})
	if err := c.Execute(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	t.Log(out.String())
}
