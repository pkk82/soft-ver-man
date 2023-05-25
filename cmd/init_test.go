package cmd_test

import (
	"github.com/pkk82/soft-ver-man/cmd"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"testing"
)

func Test_ShouldCreateConfigFile(t *testing.T) {
	tempDir := givenTempDir(t)
	execute("init")
	_, err := os.Stat(filepath.Join(tempDir, ".soft-ver-man", "config.yml"))
	if err != nil {
		t.Errorf("File was not created")
	}
}

func givenTempDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "soft-ver-man-tests-")
	if err != nil {
		t.Fatalf("Problem with creating temp dir: %v\n", err)
	}
	viper.Set(cmd.ConfigDir, tempDir)
	return tempDir
}
