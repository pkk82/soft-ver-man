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
package cmd_test

import (
	"github.com/pkk82/soft-ver-man/cmd"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_ShouldCreateConfigFile(t *testing.T) {
	viper.Reset()
	tempDir := givenTempDir(t)
	viper.Set(cmd.ConfigDir, tempDir)

	input := inputReader{inputs: []string{"/download-path\n", "/software-path\n"}, offset: 0}

	execute("init", &input)
	filePath := filepath.Join(tempDir, ".soft-ver-man", "config.yml")
	_, err := os.Stat(filePath)
	if err != nil {
		t.Errorf("File was not created")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	if !strings.Contains(string(content), "software-directory-download: /download-path") {
		t.Errorf("File %v does not contain expected content", string(content))
	}
	if !strings.Contains(string(content), "software-directory: /software-path") {
		t.Errorf("File %v does not contain expected content", string(content))
	}

}

func Test_ShouldCreateDefaultConfigFile(t *testing.T) {

	viper.Reset()
	tempDir := givenTempDir(t)
	viper.Set(cmd.ConfigDir, tempDir)
	execute("init --default", nil)

	filePath := filepath.Join(tempDir, ".soft-ver-man", "config.yml")
	_, err := os.Stat(filePath)
	if err != nil {
		t.Errorf("File was not created")
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	if !strings.Contains(string(content), "software-directory-download: /tmp/soft-ver-man") {
		t.Errorf("File %v does not contain expected content", string(content))
	}
	dir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Error reading user home dir: %v", err)
	}
	if !strings.Contains(string(content), "software-directory: "+dir+"/pf") {
		t.Errorf("File %v does not contain expected content", string(content))
	}
}

func givenTempDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "soft-ver-man-tests-")
	if err != nil {
		t.Fatalf("Problem with creating temp dir: %v\n", err)
	}
	return tempDir
}

type inputReader struct {
	inputs []string
	offset int
}

func (ir *inputReader) Read(p []byte) (n int, err error) {
	if ir.offset >= len(ir.inputs) {
		return 0, io.EOF
	}
	n = copy(p, ir.inputs[ir.offset])
	ir.offset++
	return n, nil
}
