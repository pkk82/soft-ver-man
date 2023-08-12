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
package kotlin

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/config"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/software/kotlin"
	"github.com/pkk82/soft-ver-man/util/console"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Install kotlin package into software directory",
	Long:  fmt.Sprintf("Install kotlin package (from GitHub releases) from %v into software directory", kotlin.ReleasesURL),
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		return domain.ValidateVersion(args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		softwareDownloadDir := config.InitSoftwareDownloadDir(cmd)
		softwareDir := config.InitSoftwareDir(cmd)

		fetchedPackage, err := kotlin.Fetch(args[0], softwareDownloadDir)
		if err != nil {
			console.Fatal(err)
		}
		err = kotlin.Install(fetchedPackage, softwareDir)
		if err != nil {
			console.Fatal(err)
		}
	},
}

func init() {
	Cmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}