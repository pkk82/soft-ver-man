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
package java

import (
	"github.com/pkk82/soft-ver-man/console"
	"github.com/pkk82/soft-ver-man/download"
	"github.com/pkk82/soft-ver-man/pack"
	"github.com/pkk82/soft-ver-man/verification"
	"github.com/pkk82/soft-ver-man/version"
	"path/filepath"
)

func Fetch(inputVersion, softwareDownloadDir string, verify bool) (pack.FetchedPackage, error) {
	supportedPackages := getSupportedPackages()
	versions := make([]string, len(supportedPackages))
	for i, v := range supportedPackages {
		versions[i] = v.version()
	}
	foundVersion, index, err := version.FindVersion(inputVersion, versions)
	if err != nil {
		console.Fatal(err)
	}
	matchingVersion := supportedPackages[index]
	javaDir := filepath.Join(softwareDownloadDir, Name)
	fetchedPackagePath := download.FetchFile(matchingVersion.DownloadUrl, javaDir, matchingVersion.Name)

	if verify {
		extendedPackage, err := getExtendedPackage(matchingVersion.Id)
		if err != nil {
			return pack.FetchedPackage{}, err
		}
		err = verification.VerifySha256(fetchedPackagePath, extendedPackage.Sha256)
		if err != nil {
			return pack.FetchedPackage{}, err
		}
		console.Info("Verification successful")
	}

	return pack.FetchedPackage{Version: foundVersion, FilePath: fetchedPackagePath, Type: matchingVersion.packagingType()}, nil

}
