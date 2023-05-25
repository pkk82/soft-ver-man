package cmd_test

import (
	"bytes"
	"github.com/pkk82/soft-ver-man/cmd"
	"strings"
	"testing"
)

func Test_ShouldExecuteRootCommand(t *testing.T) {
	output := execute("")
	expectedOutput := "Use \"soft-ver-man [command] --help\" for more information about a command."
	if strings.Index(output, expectedOutput) == -1 {
		t.Errorf("Root command does not contain %v\n", expectedOutput)
	}
}

func execute(args string) string {
	actual := new(bytes.Buffer)
	cmd.RootCmd.SetOut(actual)
	cmd.RootCmd.SetErr(actual)
	cmd.RootCmd.SetArgs(strings.Split(args, " "))
	err := cmd.RootCmd.Execute()
	if err != nil {
		return err.Error()
	}
	return actual.String()
}
