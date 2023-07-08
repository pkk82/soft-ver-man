package archive

import (
	"github.com/pkk82/soft-ver-man/pack"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var date = time.Now().Format("20060102-150405")

func TestExtractTarGZ(t *testing.T) {

	testDir := filepath.Join(os.TempDir(), "soft-ver-man-test-"+date, t.Name())
	err := os.MkdirAll(testDir, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create test directory: %s", err)
	}

	fetchedPackage := pack.FetchedPackage{
		FilePath: filepath.Join("testdata", "dir.tar.gz"),
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
	assertLink(t, filepath.Join(testDir, "dir", "file-11.txt"), filepath.Join("dir-1", "file-11.txt"))
}

func TestExtractZip(t *testing.T) {

	testDir := filepath.Join(os.TempDir(), "soft-ver-man-test-"+date, t.Name())
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
	assertLink(t, filepath.Join(testDir, "dir", "file-11.txt"), filepath.Join("dir-1", "file-11.txt"))

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

	if string(content) != expectedContent {
		t.Errorf("File (%s) content (%s) is not as expected: %s",
			filePath,
			string(content),
			expectedContent)
	}

}

func assertLink(t *testing.T, linkPath, filePath string) {

	isSymlink, err := isSymlink(linkPath)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	if !isSymlink {
		t.Errorf("Path (%s) is not a symlink", linkPath)
		return
	}

	target, err := getSymlinkTarget(linkPath)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	if target != filePath {
		t.Errorf("Symlink target (%s) is not as expected (%s)", target, filePath)
		return
	}

}

func isSymlink(path string) (bool, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.Mode()&os.ModeSymlink != 0, nil
}

func getSymlinkTarget(path string) (string, error) {
	targetPath, err := os.Readlink(path)
	if err != nil {
		return "", err
	}

	return targetPath, nil
}
