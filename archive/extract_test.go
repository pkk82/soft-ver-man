package archive

import (
	"github.com/pkk82/soft-ver-man/pack"
	"github.com/pkk82/soft-ver-man/test"
	"github.com/pkk82/soft-ver-man/version"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExtract(t *testing.T) {
	type args struct {
		fetchedPackage pack.FetchedPackage
	}
	tests := []struct {
		name     string
		args     args
		wantPath string
	}{
		{
			name: "tar.gz dir",
			args: args{
				fetchedPackage: pack.FetchedPackage{
					Version:  version.Version{Value: "v20.1.2"},
					FilePath: filepath.Join("testdata", "dir.tar.gz"),
					Type:     pack.TAR_GZ,
				},
			},
			wantPath: "dir",
		}, {
			name: "zip dir",
			args: args{
				fetchedPackage: pack.FetchedPackage{
					Version:  version.Version{Value: "v20.1.2"},
					FilePath: filepath.Join("testdata", "dir.zip"),
					Type:     pack.ZIP,
				},
			},
			wantPath: "dir",
		}, {
			name: "tar.gz files",
			args: args{
				fetchedPackage: pack.FetchedPackage{
					Version:  version.Version{Value: "v20.1.2"},
					FilePath: filepath.Join("testdata", "files.tar.gz"),
					Type:     pack.TAR_GZ,
				},
			},
			wantPath: "files",
		}, {
			name: "zip files",
			args: args{
				fetchedPackage: pack.FetchedPackage{
					Version:  version.Version{Value: "v20.1.2"},
					FilePath: filepath.Join("testdata", "files.zip"),
					Type:     pack.ZIP,
				},
			},
			wantPath: "files",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := test.CreateTestDir(t)
			got, err := Extract(tt.args.fetchedPackage, testDir)
			if err != nil {
				t.Errorf("Extract() error = %v", err)
				return
			}
			expectedPath := filepath.Join(testDir, tt.wantPath)
			expected := ExtractedPackage{
				Version: tt.args.fetchedPackage.Version,
				Path:    expectedPath,
			}
			if !reflect.DeepEqual(got, expected) {
				t.Errorf("Extract() got = %v, want %v", got, expected)
			}
			assertContent(t, expectedPath)
		})
	}
}

func assertContent(t *testing.T, dir string) {
	assertFileContent(t, filepath.Join(dir, "dir-1", "file-11.txt"), "file-11")
	assertFileContent(t, filepath.Join(dir, "dir-1", "file-12.txt"), "file-12")
	assertFileContent(t, filepath.Join(dir, "dir-2", "file-21.txt"), "file-21")
	assertFileContent(t, filepath.Join(dir, "dir-2", "file-22.txt"), "file-22")
	assertFileContent(t, filepath.Join(dir, "file-1.txt"), "file-1")
	assertFileContent(t, filepath.Join(dir, "exec"), "echo \"Hello, world\"\n")
	assertFileMode(t, filepath.Join(dir, "exec"), 0764)
	assertLink(t, filepath.Join(dir, "file-11.txt"), filepath.Join("dir-1", "file-11.txt"))
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

func assertFileMode(t *testing.T, filePath string, mode os.FileMode) {

	lstat, err := os.Lstat(filePath)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	if lstat.Mode() != mode {
		t.Errorf("File mode (%s) is not as expected (%s)", lstat.Mode(), mode)
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
