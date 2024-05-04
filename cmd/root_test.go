/*
 * Copyright © 2024 Piotr Kozak <piotrkrzysztofkozak@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */
package cmd_test

import (
	"bytes"
	"github.com/pkk82/soft-ver-man/cmd"
	"github.com/spf13/viper"
	"io"
	"strings"
	"testing"
)

func Test_ShouldExecuteRootCommand(t *testing.T) {
	viper.Reset()
	output := execute("", nil)
	expectedOutput := "Use \"soft-ver-man [command] --help\" for more information about a command."
	if strings.Index(output, expectedOutput) == -1 {
		t.Errorf("Root command does not contain %v\n", expectedOutput)
	}
}

func execute(args string, newIn io.Reader) string {
	actual := new(bytes.Buffer)
	cmd.RootCmd.SetOut(actual)
	cmd.RootCmd.SetErr(actual)
	cmd.RootCmd.SetArgs(strings.Split(args, " "))
	cmd.RootCmd.SetIn(newIn)
	err := cmd.RootCmd.Execute()
	if err != nil {
		return err.Error()
	}
	return actual.String()
}
