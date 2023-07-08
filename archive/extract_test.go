package archive

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/pack"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExtractTarGZ(t *testing.T) {

	testDir := filepath.Join(os.TempDir(), "soft-ver-man-test", t.Name(), fmt.Sprintf("%d", time.Now().UnixMilli()))
	err := os.MkdirAll(testDir, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create test directory: %s", err)
	}

	fetchedPackage := pack.FetchedPackage{
		FilePath: filepath.Join("testdata", "dir.tzg"),
		Type:     pack.TAR_GZ,
	}

	err = Extract(fetchedPackage, testDir)
	if err != nil {
		t.Errorf("Failed to extract tar.gz: %s", err)
	}

	assertFileContent(t, filepath.Join(testDir, "dir", "dir-1", "file-11.txt"), "file-11")
	assertFileContent(t, filepath.Join(testDir, "dir", "dir-1", "file-12.txt"), "file-12")
	assertFileContent(t, filepath.Join(testDir, "dir", "dir-2", "file-21.txt"), "file-21")
	assertFileContent(t, filepath.Join(testDir, "dir", "dir-2", "file-22.txt"), "file-22")
	assertFileContent(t, filepath.Join(testDir, "dir", "file-1.txt"), "file-1")

}

func TestExtractZip(t *testing.T) {

	testDir := filepath.Join(os.TempDir(), "soft-ver-man-test", t.Name(), fmt.Sprintf("%d", time.Now().UnixMilli()))
	err := os.MkdirAll(testDir, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create test directory: %s", err)
	}

	fetchedPackage := pack.FetchedPackage{
		FilePath: filepath.Join("testdata", "dir.zip"),
		Type:     pack.ZIP,
	}

	err = Extract(fetchedPackage, testDir)
	if err != nil {
		t.Errorf("Failed to extract zip: %s", err)
	}

	assertFileContent(t, filepath.Join(testDir, "dir", "dir-1", "file-11.txt"), "file-11")
	assertFileContent(t, filepath.Join(testDir, "dir", "dir-1", "file-12.txt"), "file-12")
	assertFileContent(t, filepath.Join(testDir, "dir", "dir-2", "file-21.txt"), "file-21")
	assertFileContent(t, filepath.Join(testDir, "dir", "dir-2", "file-22.txt"), "file-22")
	assertFileContent(t, filepath.Join(testDir, "dir", "file-1.txt"), "file-1")

}

func assertFileContent(t *testing.T, filePath, expectedContent string) {
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

	// Perform your checks on the file content
	// For example, check if it contains a specific string
	if string(content) != expectedContent {
		t.Errorf("File (%s) content (%s) is not as expected: %s", filePath, string(content), expectedContent)
	}

}
