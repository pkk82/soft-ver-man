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
package svm

import (
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/download"
	"path/filepath"
)

func Fetch(inputVersion, softwareDownloadDir string) (domain.FetchedPackage, error) {
	supportedPackages, err := getSupportedPackages()
	if err != nil {
		return domain.FetchedPackage{}, err
	}
	versions := make([]string, len(supportedPackages))
	for i, v := range supportedPackages {
		versions[i] = v.Version.Value
	}
	foundVersion, index, err := domain.FindVersion(inputVersion, versions)
	if err != nil {
		return domain.FetchedPackage{}, err
	}
	matchingVersion := supportedPackages[index]
	softDownloadDir := filepath.Join(softwareDownloadDir, Name)
	fetchedPackagePath := download.FetchFile(matchingVersion.Url, softDownloadDir, matchingVersion.Name)

	return domain.FetchedPackage{Version: foundVersion, FilePath: fetchedPackagePath, Type: matchingVersion.Type}, nil

}
