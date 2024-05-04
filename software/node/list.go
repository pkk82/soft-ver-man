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
package node

import (
	"encoding/json"
	"fmt"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/console"
	"io"
	"net/http"
	"runtime"
)

type PackagesPerVersion struct {
	Version string   `json:"version"`
	Date    string   `json:"date"`
	Files   []string `json:"files"`
}

type Package struct {
	Version  string
	FileName string
	Type     domain.Type
}

func (v Package) DownloadLink() string {
	return fmt.Sprintf("%s/%s/%s", DistURL, v.Version, v.FileName)
}
func (v Package) SumsLink() string {
	return fmt.Sprintf("%s/%s/%s", DistURL, v.Version, ShaSumFileName)
}

func (v Package) SumsSigLink() string {
	return fmt.Sprintf("%s/%s/%s", DistURL, v.Version, ShaSumSigFileName)
}

func getSupportedPackages() []Package {
	resp, err := http.Get(JsonFileURL)
	if err != nil {
		console.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			console.Fatal(err)
		}
	}(resp.Body)
	var filesPerVersions []PackagesPerVersion
	err = json.NewDecoder(resp.Body).Decode(&filesPerVersions)
	if err != nil {
		console.Fatal(err)
	}
	return supportedPackages(&filesPerVersions, runtime.GOOS, runtime.GOARCH)
}

func supportedPackages(packagesPerVersions *[]PackagesPerVersion, goOpSystem, goarch string) []Package {
	result := make([]Package, 0)
	supportedFile, supportedType := calculateSupportedFile(goOpSystem, goarch)
	for _, packagePerVersion := range *packagesPerVersions {
		if includes(packagePerVersion.Files, supportedFile) {
			fileName := calculateFileName(packagePerVersion.Version, goOpSystem, goarch)
			version := Package{
				Version:  packagePerVersion.Version,
				FileName: fileName,
				Type:     supportedType,
			}
			result = append(result, version)
		}
	}
	return result
}

// https://nodejs.org/dist/v20.3.1/node-v20.3.1-win-x64.zip
// https://nodejs.org/dist/v20.3.1/node-v20.3.1-darwin-arm64.tar.gz
// https://nodejs.org/dist/v20.3.1/node-v20.3.1-darwin-x64.tar.gz
// https://nodejs.org/dist/v20.3.1/node-v20.3.1-linux-x64.tar.gz
func calculateFileName(version, goOpSystem, goArch string) string {
	arch := toFilesArch(goArch)
	extension := toExtension(goOpSystem)
	opSys := toLinkOs(goOpSystem)
	return fmt.Sprintf("%s-%s-%s-%s.%s", Name, version, opSys, arch, extension)
}

func includes(c []string, term string) bool {
	for _, elem := range c {
		if elem == term {
			return true
		}
	}
	return false
}

func calculateSupportedFile(goOpSystem, goArch string) (string, domain.Type) {
	var supportedFile string
	var supportedType domain.Type
	os := toFilesOs(goOpSystem)
	arch := toFilesArch(goArch)
	switch os {
	case "win":
		supportedFile = fmt.Sprintf("%s-%s-zip", os, arch)
		supportedType = domain.ZIP
	case "osx":
		supportedFile = fmt.Sprintf("%s-%s-tar", os, arch)
		supportedType = domain.TAR_GZ
	default:
		supportedFile = fmt.Sprintf("%s-%s", os, arch)
		supportedType = domain.TAR_GZ
	}
	return supportedFile, supportedType
}

func toFilesOs(goOpSystem string) string {
	var os string
	switch goOpSystem {
	case "darwin":
		os = "osx"
	case "windows":
		os = "win"
	default:
		os = goOpSystem
	}
	return os
}

func toLinkOs(goOpSystem string) string {
	var os string
	switch goOpSystem {
	case "windows":
		os = "win"
	default:
		os = goOpSystem
	}
	return os
}

func toExtension(goOpSystem string) string {
	var extension string
	switch goOpSystem {
	case "windows":
		extension = "zip"
	default:
		extension = "tar.gz"
	}
	return extension
}

func toFilesArch(goArch string) string {
	var arch string
	switch goArch {
	case "amd64":
		arch = "x64"
	case "386":
		arch = "x86"
	default:
		arch = goArch
	}
	return arch
}
