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
const PackageHistorySuffix = "-domain"

const VarNameSvmSoftDir = "SVM_SOFT_DIR"
const VarNameSvmSoftPackageDirTemplate = "SVM_SOFT_%v_DIR"

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
