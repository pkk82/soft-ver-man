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
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/pkk82/soft-ver-man/console"
	"github.com/pkk82/soft-ver-man/pack"
	"io"
	"os"
	"path"
)

func Extract(fetchedPackage pack.FetchedPackage, softwareDir string) error {

	file, err := os.Open(fetchedPackage.FilePath)
	if err != nil {
		return err
	}
	defer func(open *os.File) {
		err := open.Close()
		if err != nil {
			console.Error(err)
		}
	}(file)

	if fetchedPackage.Type == pack.TAR_GZ {
		return extractTarGz(file, softwareDir)
	} else if fetchedPackage.Type == pack.ZIP {
		return extractZip(file, softwareDir)
	}
	return errors.New("Unknown archive type: " + string(fetchedPackage.Type))

}

func extractZip(file *os.File, dir string) error {
	return nil
}

func extractTarGz(file *os.File, dir string) error {
	reader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer func(reader *gzip.Reader) {
		err := reader.Close()
		if err != nil {
			console.Error(err)
		}
	}(reader)

	tarReader := tar.NewReader(reader)

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

		default:
			return errors.New(fmt.Sprintf(
				"ExtractTarGz: uknown type: %x in %s",
				header.Typeflag,
				header.Name))
		}

	}
	return nil

}
