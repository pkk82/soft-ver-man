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
package archive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/console"
	io2 "github.com/pkk82/soft-ver-man/util/io"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ExtractedPackage struct {
	Version domain.Version
	Path    string
}

func Extract(fetchedPackage domain.FetchedPackage, softwareDir string, targetDirNameStrategy domain.ExtractStrategy) (ExtractedPackage, error) {

	var err error
	var targetDirPath string
	if fetchedPackage.Type == domain.TAR_GZ {
		targetDirPath, err = extractTarGz(fetchedPackage.FilePath, softwareDir, targetDirNameStrategy)
	} else if fetchedPackage.Type == domain.ZIP {
		targetDirPath, err = extractZip(fetchedPackage.FilePath, softwareDir, targetDirNameStrategy)
	} else {
		return ExtractedPackage{}, errors.New("Unknown archive type: " + string(fetchedPackage.Type))
	}

	if err != nil {
		return ExtractedPackage{}, err
	} else {
		return ExtractedPackage{Version: fetchedPackage.Version, Path: targetDirPath}, nil
	}
}

func extractZip(zipPath string, dir string, strategy domain.ExtractStrategy) (string, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer io2.CloseOrLog(reader)

	topLevelDir := extractTopLevelDirInZipFile(reader)
	targetFilePathSupplier := prepareTargetFilePathSupplier(zipPath, topLevelDir, strategy)

	for _, file := range reader.File {
		targetFilePath := targetFilePathSupplier.supply(dir, file.Name)
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(targetFilePath, file.Mode())
			if err != nil {
				return "", err
			}
			continue
		}
		err = extractZipFile(targetFilePath, file)
		if err != nil {
			return "", err
		}
	}

	return targetFilePathSupplier.supply(dir, ""), nil
}

func extractTopLevelDirInZipFile(reader *zip.ReadCloser) string {
	var topLevelDirs []string
	for _, file := range reader.File {
		if !strings.Contains(strings.TrimSuffix(file.Name, "/"), "/") && file.FileInfo().IsDir() {
			topLevelDirs = append(topLevelDirs, file.Name)
		}
	}
	if len(topLevelDirs) == 1 {
		return topLevelDirs[0]
	} else {
		return ""
	}
}

func extractZipFile(targetFilePath string, file *zip.File) error {

	zipEntry, err := file.Open()
	if err != nil {
		return err
	}
	defer io2.CloseOrLog(zipEntry)

	if file.Mode()&os.ModeSymlink != 0 {
		link, err := io.ReadAll(zipEntry)

		if err != nil {
			return err
		}

		if err := os.Symlink(string(link), targetFilePath); err != nil {
			return err
		}
		return nil
	}

	targetFile, err := os.OpenFile(targetFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer io2.CloseOrLog(targetFile)

	_, err = io.Copy(targetFile, zipEntry)
	if err != nil {
		return err
	}

	fileMode := file.Mode()
	if fileMode&0111 != 0 {
		err := os.Chmod(targetFilePath, fileMode|0100)
		if err != nil {
			return errors.New(fmt.Sprintf("ExtractTarGz: Chmod() failed: %s", err.Error()))
		}
	}

	return nil
}

func extractTarGz(tarGzFilePath, dir string, strategy domain.ExtractStrategy) (string, error) {

	topLevelDir, err := extractTopLevelDirInTarGzFile(tarGzFilePath)
	if err != nil {
		return "", err
	}
	targetFilePathSupplier := prepareTargetFilePathSupplier(tarGzFilePath, topLevelDir, strategy)

	tarGzFile, err := os.Open(tarGzFilePath)
	if err != nil {
		return "", err
	}
	defer io2.CloseOrLog(tarGzFile)

	gzReader, err := gzip.NewReader(tarGzFile)
	if err != nil {
		return "", err
	}
	defer io2.CloseOrLog(gzReader)

	tarReader := tar.NewReader(gzReader)
	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			return targetFilePathSupplier.supply(dir, ""), nil
		}

		if err != nil {
			return "", err
		}

		targetFilePath := targetFilePathSupplier.supply(dir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetFilePath, 0755); err != nil {
				return "", errors.New(fmt.Sprintf("ExtractTarGz: MkdirAll() failed: %s", err.Error()))
			}
		case tar.TypeReg:
			outFile, err := os.Create(targetFilePath)
			if err != nil {
				return "", errors.New(fmt.Sprintf("ExtractTarGz: Create() failed: %s", err.Error()))
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return "", errors.New(fmt.Sprintf("ExtractTarGz: Copy() failed: %s", err.Error()))
			}
			err = outFile.Close()
			if err != nil {
				return "", err
			}

			fileMode := header.FileInfo().Mode()
			if fileMode&0111 != 0 {
				err := os.Chmod(targetFilePath, fileMode|0100)
				if err != nil {
					return "", errors.New(fmt.Sprintf("ExtractTarGz: Chmod() failed: %s", err.Error()))
				}
			}

		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, targetFilePath); err != nil {
				console.Fatal(err)
				return "", err
			}

		default:
			return "", errors.New(fmt.Sprintf(
				"ExtractTarGz: uknown type: %x in %s",
				header.Typeflag,
				header.Name))
		}

	}
	return "", errors.New("should not reach here")

}

func extractTopLevelDirInTarGzFile(tarGzFilePath string) (string, error) {

	tarGzFile, err := os.Open(tarGzFilePath)
	if err != nil {
		return "", err
	}
	defer io2.CloseOrLog(tarGzFile)

	gzReader, err := gzip.NewReader(tarGzFile)
	if err != nil {
		return "", err
	}
	defer io2.CloseOrLog(gzReader)

	tarReader := tar.NewReader(gzReader)

	topLevelDirs := make([]string, 0)
	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			if len(topLevelDirs) == 1 {
				return topLevelDirs[0], nil
			} else {
				return "", nil
			}
		}

		if err != nil {
			return "", err
		}

		if header.Typeflag == tar.TypeDir {
			if !strings.Contains(strings.TrimSuffix(header.Name, "/"), "/") && !contains(topLevelDirs, header.Name) {
				topLevelDirs = append(topLevelDirs, header.Name)
			}
		}
	}
	return "", errors.New("should not reach here")
}

func archiveNameWithoutExtension(path string) string {
	var name = filepath.Base(path)
	if strings.HasSuffix(name, "."+domain.TAR_GZ) {
		return strings.TrimSuffix(name, "."+domain.TAR_GZ)
	} else {
		return strings.TrimSuffix(name, filepath.Ext(name))
	}
}

type targetFilePathSupplier interface {
	supply(targetDir, archiveFilePath string) string
}

type defaultTargetFilePathSupplier struct {
	archiveNameWithoutExtension string
	topLevelDirInArchive        string
}

func (s *defaultTargetFilePathSupplier) supply(targetDir, archiveFilePath string) string {
	if s.topLevelDirInArchive == "" {
		return filepath.Join(targetDir, s.archiveNameWithoutExtension, archiveFilePath)
	} else {
		return filepath.Join(targetDir, s.topLevelDirInArchive, strings.TrimPrefix(archiveFilePath, s.topLevelDirInArchive))
	}
}

type archiveReplaceTargetFilePathSupplier struct {
	archiveNameWithoutExtension string
	topLevelDirInArchive        string
}

func (s *archiveReplaceTargetFilePathSupplier) supply(targetDir, archiveFilePath string) string {
	if s.topLevelDirInArchive == "" {
		return filepath.Join(targetDir, s.archiveNameWithoutExtension, archiveFilePath)
	} else {
		return filepath.Join(targetDir, s.archiveNameWithoutExtension, strings.TrimPrefix(archiveFilePath, s.topLevelDirInArchive))
	}
}

func prepareTargetFilePathSupplier(archivePath, topLevelDirInArchive string, strategy domain.ExtractStrategy) targetFilePathSupplier {

	archiveNameWithoutExtension := archiveNameWithoutExtension(archivePath)
	switch strategy {
	case domain.UseCompressedDirOrArchiveName:
		return &defaultTargetFilePathSupplier{archiveNameWithoutExtension: archiveNameWithoutExtension, topLevelDirInArchive: topLevelDirInArchive}
	case domain.ReplaceCompressedDirWithArchiveName:
		return &archiveReplaceTargetFilePathSupplier{archiveNameWithoutExtension: archiveNameWithoutExtension, topLevelDirInArchive: topLevelDirInArchive}
	default:
		panic("Unknown target dir name strategy: " + string(strategy))
	}

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
