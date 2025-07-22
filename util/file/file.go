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

package file

import (
	"errors"
	"github.com/pkk82/soft-ver-man/domain"
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

func Extension(fp string) domain.Type {
	ext := filepath.Ext(fp)
	name := filepath.Base(fp)
	if ext == ".zip" {
		return domain.ZIP
	} else if ext == ".gz" && strings.HasSuffix(name, ".tar.gz") {
		return domain.TAR_GZ
	} else {
		return domain.UNKNOWN
	}
}

func createFile(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			console.Error(err)
		}
	}(f)

	return nil
}

func AssertFileWithContent(filePath string, seekLine string, contentToAdd []string) error {

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

		err = createFile(filePath)
		if err != nil {
			return err
		}
	}

	content, err := ReadFile(filePath)
	if err != nil {
		return err
	}

	if !strings.Contains(content, seekLine) {
		err := AppendInFile(filePath, contentToAdd)
		if err != nil {
			return err
		}

	}

	return nil

}
