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
package cmd

import (
	"bufio"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const SoftwareDirKey = "software-directory"

func InitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Init soft-ver-man",
		Long:  `Init soft-ver-man by defining directory where software is to be installed and install itself there`,
		Run: func(cmd *cobra.Command, args []string) {
			if !viper.IsSet(SoftwareDirKey) {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					homeDir = ""
				}
				if homeDir == "" {
					cmd.Printf("Where to install software: ")
				} else {
					cmd.Printf("Where to install software [%v]: ", filepath.Join(homeDir, "pf"))
				}
				reader := bufio.NewReader(cmd.InOrStdin())
				softwareDirInput, err := reader.ReadString('\n')
				if err != nil {
					cmd.PrintErr(err)
				}
				if strings.TrimSpace(softwareDirInput) == "" && homeDir != "" {
					viper.Set(SoftwareDirKey, filepath.Join(homeDir, "pf"))
				} else if strings.TrimSpace(softwareDirInput) != "" {
					viper.Set(SoftwareDirKey, softwareDirInput)
				} else {
					cobra.CheckErr("Software directory not set")
				}
				err = viper.WriteConfig()
				cobra.CheckErr(err)
			}
		},
	}
}

func init() {
	RootCmd.AddCommand(InitCommand())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}