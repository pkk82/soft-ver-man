package config

import (
	"bufio"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const SoftwareDownloadDirKey = "software-directory-download"
const SoftwareDirKey = "software-directory"

func InitSoftwareDownloadDir(cmd *cobra.Command) string {
	return initConfigEntry(cmd,
		SoftwareDownloadDirKey,
		"Where to download software?",
		"Software download directory not set",
		func() (string, error) {
			return filepath.Join(os.TempDir(), "soft-ver-man"), nil
		})

}

func InitSoftwareDir(cmd *cobra.Command) string {
	return initConfigEntry(cmd,
		SoftwareDirKey,
		"Where to install software?",
		"Software directory not set",
		func() (string, error) {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			return filepath.Join(homeDir, "pf"), nil
		})
}

func initConfigEntry(cmd *cobra.Command,
	key, question, errMessage string,
	defaultValueSupplier func() (string, error)) string {
	if !viper.IsSet(key) {
		defaultValue, err := defaultValueSupplier()
		if err != nil {
			defaultValue = ""
		}
		if defaultValue == "" {
			cmd.Printf(question + ": ")
		} else {
			cmd.Printf("%v [%v]: ", question, defaultValue)
		}
		reader := bufio.NewReader(cmd.InOrStdin())
		input, err := reader.ReadString('\n')
		if err != nil {
			cmd.PrintErr(err)
		}
		if strings.TrimSpace(input) == "" && defaultValue != "" {
			viper.Set(key, defaultValue)
		} else if strings.TrimSpace(input) != "" {
			viper.Set(key, input)
		} else {
			cobra.CheckErr(errMessage)
		}
		err = viper.WriteConfig()
		cobra.CheckErr(err)
	}
	return viper.GetString(key)
}
