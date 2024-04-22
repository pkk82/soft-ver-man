package file_test

import (
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
