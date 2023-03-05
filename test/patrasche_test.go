package test

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"

	"github.com/key-inside/patrasche"
	"github.com/key-inside/patrasche/cmd/ccquery"
	"github.com/key-inside/patrasche/cmd/inspect"
)

const testConfigs = "--config=./fixtures/config.yaml"

func bufferdOutput(c *cobra.Command) (out *bytes.Buffer) {
	out = bytes.NewBuffer(nil)
	c.SetOut(out)
	c.SetErr(out)
	return
}

func Test_Help(t *testing.T) {
	c := &cobra.Command{
		Use:     "test",
		Short:   "Patrasche Test",
		Long:    "Patrasche Test Command",
		Version: "v0.0.0",
	}
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

func Test_Inspect(t *testing.T) {
	c := &cobra.Command{
		Use:     "test",
		Short:   "Patrasche Test",
		Long:    "Patrasche Test Command",
		Version: "v0.0.0",
	}
	c.AddCommand(inspect.Command())
	// sets output before biting
	out := bufferdOutput(c)

	_, err := patrasche.Bite(c, patrasche.WithConsoleLogWriter())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

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
	c := &cobra.Command{
		Use:     "test",
		Short:   "Patrasche Test",
		Long:    "Patrasche Test Command",
		Version: "v0.0.0",
	}
	c.AddCommand(inspect.Command(), ccquery.Command())
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
