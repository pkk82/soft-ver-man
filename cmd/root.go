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
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "soft-ver-man",
	Short: "Install supported software in a transparent way and easily switch between versions",
	Long: `Install supported software in a transparent way and easily switch between different versions of it.

The app can download or copy supported software in many versions to specified folder on disk.
Operating only with environment variables it can switch between different versions of given software.

Additional documentation can be found in https://github.com/pkk82/soft-ver-man

To start just run: soft-ver-man init`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

const ConfigDir = "config-directory"

func init() {
	cobra.OnInitialize(initConfig)
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// find config directory
	configDir := viper.GetString(ConfigDir)
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)
		configDir = homeDir
	}

	configPath := filepath.Join(configDir, ".soft-ver-man")
	err := os.Mkdir(configPath, 0700)
	if err != nil && !os.IsExist(err) {
		displayError(err)
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigType("yml")
	viper.SetConfigName("config")

	// read in environment variables that match
	viper.AutomaticEnv()

	err = viper.SafeWriteConfig()
	// SaveWriteConfig return error that could not be checked with os.IsExist
	_, statErr := os.Stat(filepath.Join(configPath, "config.yml"))
	if err != nil && statErr != nil {
		displayError(err)
	}

	if err := viper.ReadInConfig(); err == nil {
		displayMessageOnStdErr("Using config file:", viper.ConfigFileUsed())
	}
}

func displayMessageOnStdErr(msgs ...any) {
	_, err := fmt.Fprintln(os.Stderr, msgs...)
	if err != nil {
		os.Exit(1)
	}
}

func displayError(error error) {
	displayMessageOnStdErr("Error:", error.Error())
}
