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

package config

import (
	"bufio"
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const HomeConfigDir = ".soft-ver-man"
const RcFile = ".svmmainrc"
const SoftwareDownloadDirKey = "software-directory-download"
const SoftwareDirKey = "software-directory"
const InstalledPackagesSuffix = "-installed-packages"

type Config struct {
	SoftwareDownloadDir string
	SoftwareDir         string
}

func Get() (Config, error) {
	if !viper.IsSet(SoftwareDownloadDirKey) || !viper.IsSet(SoftwareDirKey) {
		return Config{}, errors.New("config not initialized, run 'svm init' first")
	}
	return Config{
		SoftwareDownloadDir: viper.GetString(SoftwareDownloadDirKey),
		SoftwareDir:         viper.GetString(SoftwareDirKey),
	}, nil
}

func Init(cmd *cobra.Command, useDefault bool) {
	initSoftwareDownloadDir(cmd, useDefault)
	initSoftwareDir(cmd, useDefault)
}
func initSoftwareDownloadDir(cmd *cobra.Command, useDefault bool) {

	initConfigEntry(cmd,
		SoftwareDownloadDirKey,
		"Where to download software?",
		"Software download directory not set",
		func() (string, error) {
			return filepath.Join(os.TempDir(), "soft-ver-man"), nil
		}, useDefault)
}

func initSoftwareDir(cmd *cobra.Command, useDefault bool) {
	initConfigEntry(cmd,
		SoftwareDirKey,
		"Where to install software?",
		"Software directory not set",
		func() (string, error) {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			return filepath.Join(homeDir, "pf"), nil
		}, useDefault)
}

func initConfigEntry(cmd *cobra.Command,
	key, question, errMessage string,
	defaultValueSupplier func() (string, error),
	useDefault bool) {
	if !viper.IsSet(key) {
		defaultValue, err := defaultValueSupplier()
		if err != nil {
			defaultValue = ""
		}

		if defaultValue == "" {
			cmd.Printf(question + ": ")
		} else if !useDefault {
			cmd.Printf("%v [%v]: ", question, defaultValue)
		}
		var input string
		if useDefault && defaultValue != "" {
			input = defaultValue
		} else {
			reader := bufio.NewReader(cmd.InOrStdin())
			input, err = reader.ReadString('\n')
			if err != nil {
				cmd.PrintErr(err)
			}
		}

		if strings.TrimSpace(input) == "" && defaultValue != "" {
			viper.Set(key, defaultValue)
		} else if strings.TrimSpace(input) != "" {
			viper.Set(key, strings.TrimSpace(input))
		} else {
			cobra.CheckErr(errMessage)
		}
		err = viper.WriteConfig()
		cobra.CheckErr(err)
	}
}
