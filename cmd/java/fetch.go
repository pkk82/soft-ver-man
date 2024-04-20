/*
Copyright © 2023 Piotr Kozak <piotrkrzysztofkozak@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package java

import (
	"fmt"
	cmdUtil "github.com/pkk82/soft-ver-man/cmd"
	"github.com/pkk82/soft-ver-man/config"
	"github.com/pkk82/soft-ver-man/software/java"
	"github.com/pkk82/soft-ver-man/util/console"
	"github.com/spf13/cobra"
)

var fetchVerify bool

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:     "fetch [version]",
	Aliases: []string{"f", "fetch"},
	Short:   "Fetch java package into download directory",
	Long:    fmt.Sprintf("Fetch java package from %v into download directory", java.PackagesAPIURL),
	Args:    cmdUtil.VersionArg,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.Get()
		if err != nil {
			console.Fatal(err)
		}
		fetch, err := java.Fetch(cmdUtil.FirstOrEmpty(args), config.SoftwareDownloadDir, fetchVerify)
		if err != nil {
			console.Fatal(err)
		} else {
			console.Info(fmt.Sprintf("Downloaded file: %v", fetch.FilePath))
		}
	},
}

func init() {
	Cmd.AddCommand(fetchCmd)
	fetchCmd.Flags().BoolVarP(&fetchVerify, "verify", "", false, "Verify checksum of downloaded file")
}
