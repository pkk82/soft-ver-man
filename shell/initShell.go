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
package shell

import (
	"errors"
	"github.com/pkk82/soft-ver-man/config"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/console"
	"github.com/pkk82/soft-ver-man/util/file"
	"os"
	"path/filepath"
	"strings"
)

func initShell(finder domain.DirFinder) error {
	header := "### soft-ver-man"
	initLine := bashToLoad(config.RcFile)
	dir, err := finder.HomeDir()
	if err != nil {
		return err
	}

	atLeastOneExists := false

	rcFiles := [2]string{".bashrc", ".zshrc"}
	for _, rcFile := range rcFiles {
		rcPath := filepath.Join(dir, rcFile)
		exists, err := file.FileExists(rcPath)
		if err != nil {
			return err
		}
		if exists {
			atLeastOneExists = true
			content, err := file.ReadFile(rcPath)
			if err != nil {
				return err
			}
			if strings.HasSuffix(content, "\n") {
				header = "\n" + header
			} else if len(content) > 0 {
				header = "\n\n" + header
			}

			err = assertFileWithContent(rcPath, initLine, []string{header, initLine})
			if err != nil {
				return err
			}
		}
	}

	if !atLeastOneExists {
		return errors.New("at least of of the following files must exist: " + strings.Join(rcFiles[:], ", "))
	}
	return nil
}

func assertFileWithContent(filePath string, seekLine string, contentToAdd []string) error {

	exists, err := file.FileExists(filePath)
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

	content, err := file.ReadFile(filePath)
	if err != nil {
		return err
	}

	if !strings.Contains(content, seekLine) {
		err := file.AppendInFile(filePath, contentToAdd)
		if err != nil {
			return err
		}

	}

	return nil

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
