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

package archive

import (
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/test"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExtract(t *testing.T) {
	type args struct {
		fetchedPackage domain.FetchedPackage
		strategy       domain.ExtractStrategy
	}
	tests := []struct {
		name     string
		args     args
		wantPath string
	}{
		{
			name: "tar.gz dir (default strategy)",
			args: args{
				fetchedPackage: domain.FetchedPackage{
					Version:  domain.Version{Value: "v20.1.2"},
					FilePath: filepath.Join("testdata", "some-dir.tar.gz"),
					Type:     domain.TAR_GZ,
				},
				strategy: domain.UseCompressedDirOrArchiveName,
			},
			wantPath: "dir",
		}, {
			name: "tar.gz dir (replace strategy)",
			args: args{
				fetchedPackage: domain.FetchedPackage{
					Version:  domain.Version{Value: "v20.1.2"},
					FilePath: filepath.Join("testdata", "some-dir.tar.gz"),
					Type:     domain.TAR_GZ,
				},
				strategy: domain.ReplaceCompressedDirWithArchiveName,
			},
			wantPath: "some-dir",
		}, {
			name: "zip dir (default strategy)",
			args: args{
				fetchedPackage: domain.FetchedPackage{
					Version:  domain.Version{Value: "v20.1.2"},
					FilePath: filepath.Join("testdata", "some-dir.zip"),
					Type:     domain.ZIP,
				},
				strategy: domain.UseCompressedDirOrArchiveName,
			},
			wantPath: "dir",
		}, {
			name: "zip dir (replace strategy)",
			args: args{
				fetchedPackage: domain.FetchedPackage{
					Version:  domain.Version{Value: "v20.1.2"},
					FilePath: filepath.Join("testdata", "some-dir.zip"),
					Type:     domain.ZIP,
				},
				strategy: domain.ReplaceCompressedDirWithArchiveName,
			},
			wantPath: "some-dir",
		}, {
			name: "tar.gz files",
			args: args{
				fetchedPackage: domain.FetchedPackage{
					Version:  domain.Version{Value: "v20.1.2"},
					FilePath: filepath.Join("testdata", "files.tar.gz"),
					Type:     domain.TAR_GZ,
				},
				strategy: domain.UseCompressedDirOrArchiveName,
			},

			wantPath: "files",
		}, {
			name: "zip files",
			args: args{
				fetchedPackage: domain.FetchedPackage{
					Version:  domain.Version{Value: "v20.1.2"},
					FilePath: filepath.Join("testdata", "files.zip"),
					Type:     domain.ZIP,
				},
				strategy: domain.UseCompressedDirOrArchiveName,
			},
			wantPath: "files",
		},
		{
			name: "tar.gz dir (duplicated dir)",
			args: args{
				fetchedPackage: domain.FetchedPackage{
					Version:  domain.Version{Value: "v20.1.2"},
					FilePath: filepath.Join("testdata", "duplicated-dir.tar.gz"),
					Type:     domain.TAR_GZ,
				},
				strategy: domain.UseCompressedDirOrArchiveName,
			},
			wantPath: "dir",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := test.CreateTestDir(t)
			got, err := Extract(tt.args.fetchedPackage, testDir, tt.args.strategy)
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
