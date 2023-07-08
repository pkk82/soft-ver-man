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
	"github.com/pkk82/soft-ver-man/console"
	"github.com/pkk82/soft-ver-man/pack"
	"io"
	"os"
	"path"
	"path/filepath"
)

func Extract(fetchedPackage pack.FetchedPackage, softwareDir string) error {

	if fetchedPackage.Type == pack.TAR_GZ {
		return extractTarGz(fetchedPackage.FilePath, softwareDir)
	} else if fetchedPackage.Type == pack.ZIP {
		return extractZip(fetchedPackage.FilePath, softwareDir)
	}
	return errors.New("Unknown archive type: " + string(fetchedPackage.Type))

}

func extractZip(zipPath string, dir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer func(reader *zip.ReadCloser) {
		err := reader.Close()
		if err != nil {
			console.Error(err)
		}
	}(reader)

	for _, file := range reader.File {
		targetFilePath := filepath.Join(dir, file.Name)
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(targetFilePath, file.Mode())
			if err != nil {
				return err
			}
			continue
		}
		err = extractZipFile(targetFilePath, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractZipFile(targetFilePath string, file *zip.File) error {

	zipEntry, err := file.Open()
	if err != nil {
		return err
	}
	defer func(zipEntry io.ReadCloser) {
		err := zipEntry.Close()
		if err != nil {
			console.Error(err)
		}
	}(zipEntry)

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
	defer func(targetFile *os.File) {
		err := targetFile.Close()
		if err != nil {
			console.Error(err)
		}
	}(targetFile)

	_, err = io.Copy(targetFile, zipEntry)
	if err != nil {
		return err
	}
	return nil
}

func extractTarGz(tarGzFilePath, dir string) error {

	tarGzFile, err := os.Open(tarGzFilePath)
	if err != nil {
		return err
	}
	defer func(open *os.File) {
		err := open.Close()
		if err != nil {
			console.Error(err)
		}
	}(tarGzFile)

	gzReader, err := gzip.NewReader(tarGzFile)
	if err != nil {
		return err
	}
	defer func(reader *gzip.Reader) {
		err := reader.Close()
		if err != nil {
			console.Error(err)
		}
	}(gzReader)

	tarReader := tar.NewReader(gzReader)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path.Join(dir, header.Name), 0755); err != nil {
				return errors.New(fmt.Sprintf("ExtractTarGz: MkdirAll() failed: %s", err.Error()))
			}
		case tar.TypeReg:
			outFile, err := os.Create(path.Join(dir, header.Name))
			if err != nil {
				return errors.New(fmt.Sprintf("ExtractTarGz: Create() failed: %s", err.Error()))
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return errors.New(fmt.Sprintf("ExtractTarGz: Copy() failed: %s", err.Error()))
			}
			err = outFile.Close()
			if err != nil {
				return err
			}
		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, filepath.Join(dir, header.Name)); err != nil {
				console.Fatal(err)
				return err
			}

		default:
			return errors.New(fmt.Sprintf(
				"ExtractTarGz: uknown type: %x in %s",
				header.Typeflag,
				header.Name))
		}

	}
	return nil

}
