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
package shell

import (
	"errors"
	"github.com/pkk82/soft-ver-man/console"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func initBash(finder homeDirFinder) error {
	header := "### soft-ver-man"
	initLine := "[[ -s \"$HOME/.soft-ver-man/.svmrc\" ]] && source \"$HOME/.soft-ver-man/.svmrc\""
	dir, err := finder.HomeDir()
	if err != nil {
		return err
	}

	bashRcPath := filepath.Join(dir, ".bashrc")
	exists, err := fileExists(bashRcPath)
	if err != nil {
		return err
	}

	if exists {
		content, err := readFile(bashRcPath)
		if err != nil {
			return err
		}
		if strings.HasSuffix(content, "\n") {
			header = "\n" + header
		} else if len(content) > 0 {
			header = "\n\n" + header
		}
	}

	err = assertFileWithContent(bashRcPath, initLine, []string{header, initLine})
	if err != nil {
		return err
	}
	return nil
}

func assertFileWithContent(filePath string, seekLine string, contentToAdd []string) error {

	exists, err := fileExists(filePath)
	if err != nil {
		return err
	}
	if !exists {
		err := createFile(filePath)
		if err != nil {
			return err
		}
	}

	content, err := readFile(filePath)
	if err != nil {
		return err
	}

	if !strings.Contains(content, seekLine) {
		err := appendInFile(filePath, contentToAdd)
		if err != nil {
			return err
		}

	}

	return nil

}

func readFile(path string) (string, error) {
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

func appendInFile(path string, lines []string) error {
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

func fileExists(filePath string) (bool, error) {
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

func createFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			console.Error(err)
		}
	}(file)

	return nil
}
