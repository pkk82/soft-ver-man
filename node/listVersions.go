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
	"io"
	"log"
	"net/http"
	"runtime"
)

type FilesPerVersion struct {
	Version string   `json:"version"`
	Date    string   `json:"date"`
	Files   []string `json:"files"`
}

type Version struct {
	Id           string
	FileName     string
	DownloadLink string
}

func ListVersions() {
	versions := getSupportedVersions()
	for _, version := range versions {
		println(version.Id)
	}
}

func getSupportedVersions() []Version {
	resp, err := http.Get(JsonFileURL)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	var filesPerVersions []FilesPerVersion
	err = json.NewDecoder(resp.Body).Decode(&filesPerVersions)
	if err != nil {
		log.Fatal(err)
	}
	return supportedVersions(&filesPerVersions, runtime.GOOS, runtime.GOARCH)
}

func supportedVersions(filesPerVersions *[]FilesPerVersion, goOpSystem, goarch string) []Version {
	result := make([]Version, 0)
	expectedFile := supportedFile(goOpSystem, goarch)
	for _, filesPerVersion := range *filesPerVersions {
		if includes(filesPerVersion.Files, expectedFile) {
			fileName := calculateFileName(filesPerVersion.Version, goOpSystem, goarch)
			version := Version{
				Id:           filesPerVersion.Version,
				FileName:     fileName,
				DownloadLink: fmt.Sprintf("%s/%s/%s", DistURL, filesPerVersion.Version, fileName),
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
	return fmt.Sprintf("node-%s-%s-%s.%s", version, opSys, arch, extension)
}

func includes(c []string, term string) bool {
	for _, elem := range c {
		if elem == term {
			return true
		}
	}
	return false
}

func supportedFile(goOpSystem, goArch string) string {
	var supportedFile string
	os := toFilesOs(goOpSystem)
	arch := toFilesArch(goArch)
	switch os {
	case "win":
		supportedFile = fmt.Sprintf("%s-%s-zip", os, arch)
	case "osx":
		supportedFile = fmt.Sprintf("%s-%s-tar", os, arch)
	default:
		supportedFile = fmt.Sprintf("%s-%s", os, arch)
	}
	return supportedFile
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
