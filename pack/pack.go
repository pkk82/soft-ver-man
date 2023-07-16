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
package pack

import (
	"github.com/pkk82/soft-ver-man/version"
	"path/filepath"
	"strings"
)

const (
	TAR_GZ = "tar.gz"
	ZIP    = "zip"
)

type Type string

type FetchedPackage struct {
	Version  version.Version
	FilePath string
	Type     Type
}

func (fp FetchedPackage) getDirName() string {
	base := filepath.Base(fp.FilePath)
	if fp.Type == TAR_GZ && strings.HasSuffix(base, "."+TAR_GZ) {
		return strings.TrimSuffix(base, ".tar.gz")
	}
	return strings.TrimSuffix(base, filepath.Ext(base))
}

type InstalledPackage struct {
	Version     version.Version
	Path        string
	InstalledOn int64
}
