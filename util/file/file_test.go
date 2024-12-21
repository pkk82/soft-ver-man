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

package file_test

import (
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/file"
	"github.com/pkk82/soft-ver-man/util/test"
	"path/filepath"
	"strings"
	"testing"
)

func TestOverrideFileWithContent(t *testing.T) {

	tests := []struct {
		name            string
		existingContent []string
		newContent      []string
		expectedContent []string
	}{
		{"no file", nil, []string{"new content"}, []string{"new content"}},
		{"existing file", []string{"existing content"}, []string{"new content"}, []string{"new content"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := test.CreateTestDir(t)

			if tt.existingContent != nil {
				test.CreateFile(dir, "file-"+tt.name, tt.existingContent, t)
			}

			filePath := filepath.Join(dir, "file-"+tt.name)
			err := file.OverrideFileWithContent(filePath, tt.newContent)

			if err != nil {
				t.Errorf("Error: %s", err)
			}

			actualContent, err := file.ReadFile(filePath)

			if err != nil {
				t.Errorf("Error: %s", err)
			}

			if actualContent != strings.Join(tt.expectedContent, "\n") {
				t.Errorf("Expected content: %s, got: %s", strings.Join(tt.expectedContent, "\n"), actualContent)
			}

		})
	}
}

func TestAppendFileWithContent(t *testing.T) {

	tests := []struct {
		name            string
		existingContent []string
		newContent      []string
		expectedContent []string
	}{
		{"existing file", []string{"existing content", ""}, []string{"new content"}, []string{"existing content", "new content", ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := test.CreateTestDir(t)

			if tt.existingContent != nil {
				test.CreateFile(dir, "file-"+tt.name, tt.existingContent, t)
			}

			filePath := filepath.Join(dir, "file-"+tt.name)
			err := file.AppendInFile(filePath, tt.newContent)

			if err != nil {
				t.Errorf("Error: %s", err)
			}

			actualContent, err := file.ReadFile(filePath)

			if err != nil {
				t.Errorf("Error: %s", err)
			}

			if actualContent != strings.Join(tt.expectedContent, "\n") {
				t.Errorf("Expected content: %s, got: %s", strings.Join(tt.expectedContent, "\n"), actualContent)
			}

		})
	}
}

func TestExtension(t *testing.T) {

	tests := []struct {
		name         string
		filename     string
		expectedType domain.Type
	}{
		{"zip", "archive.zip", domain.ZIP}, {"tar.gz", "archive.tar.gz", domain.TAR_GZ},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := test.CreateTestDir(t)
			test.CreateFile(dir, tt.filename, []string{}, t)

			filePath := filepath.Join(dir, tt.filename)

			actualType := file.Extension(filePath)

			if actualType != tt.expectedType {
				t.Errorf("Expected type: %s, got: %s", tt.expectedType, actualType)
			}

		})
	}
}
