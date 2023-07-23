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
package test

import (
	"io"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func CreateTestDir(t *testing.T) string {

	date := time.Now().Format("20060102-150405") + "-" + generateString(5)
	testDir := filepath.Join(os.TempDir(), "soft-ver-man-test-"+date, t.Name())
	err := os.MkdirAll(testDir, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create test directory: %s", err)
	}
	return testDir
}

func CreateEmptyFile(path, fileName string, t *testing.T) string {
	file, err := os.Create(filepath.Join(path, fileName))
	if err != nil {
		t.Fatalf("Failed to create file: %s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			t.Errorf("Failed to close file: %s", err)
		}
	}(file)
	return file.Name()
}

func CreateFile(path, fileName string, content []string, t *testing.T) {
	err := os.WriteFile(filepath.Join(path, fileName), []byte(join(content)), os.ModePerm)
	if err != nil {
		t.Errorf("Failed to create file: %s", err)
	}
}

func join(content []string) string {
	return strings.Join(content, "\n")
}

type TestDirs struct {
	Home string
}

func (receiver TestDirs) HomeDir() (string, error) {
	return receiver.Home, nil
}

func (receiver TestDirs) SoftDir() (string, error) {
	return path.Join(receiver.Home, "pf"), nil
}

func AssertFileContent(path, fileName string, expectedContent []string, t *testing.T) {
	filePath := filepath.Join(path, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open file: %s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			t.Errorf("Failed to close file: %s", err)
		}
	}(file)

	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("Failed to read file content: %s", err)
	}

	if string(content) != join(expectedContent) {
		t.Errorf("File (%s) content:\n\n%s\n\n is not as expected:\n\n%s\n\n",
			filePath,
			string(content),
			join(expectedContent))
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateString(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, length)
	charsetLength := len(charset)

	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(charsetLength)]
	}

	return string(result)
}
