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
package copy

import (
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/console"
	"io"
	"os"
	"path"
)

type CopiedPackage struct {
	Version    domain.Version
	PathToFile string
	FileName   string
}

func Copy(fetchedPackage domain.FetchedPackage, dir, name string) (CopiedPackage, error) {

	srcFile, err := os.Open(fetchedPackage.FilePath)
	if err != nil {
		return CopiedPackage{}, err
	}
	defer func(srcFile *os.File) {
		err := srcFile.Close()
		if err != nil {
			console.Error(err)
		}
	}(srcFile)

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return CopiedPackage{}, err
	}

	dstFile, err := os.Create(path.Join(dir, name))
	if err != nil {
		return CopiedPackage{}, err
	}
	defer func(dstFile *os.File) {
		err := dstFile.Close()
		if err != nil {
			console.Error(err)
		}
	}(dstFile)

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return CopiedPackage{}, err
	}

	err = os.Chmod(dstFile.Name(), 0755)
	if err != nil {
		return CopiedPackage{}, err
	}

	return CopiedPackage{Version: fetchedPackage.Version, PathToFile: dir, FileName: name}, nil
}
