package file

import (
	"errors"
	"github.com/pkk82/soft-ver-man/util/console"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func OverrideFileWithContent(filePath string, contentToAdd []string) error {

	exists, err := FileExists(filePath)
	if err != nil {
		return err
	}
	if !exists {
		parent, _ := filepath.Split(filePath)
		err := os.MkdirAll(parent, 0755)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(filePath, []byte(strings.Join(contentToAdd, "\n")), 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			console.Error(err)
		}
	}(file)

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func AppendInFile(path string, lines []string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			console.Error(err)
		}
	}(file)

	for _, line := range lines {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func FileExists(filePath string) (bool, error) {
	fileinfo, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	mode := fileinfo.Mode()
	if !mode.IsRegular() {
		return false, errors.New("Not a regular file: " + filePath)
	}

	return true, nil
}
