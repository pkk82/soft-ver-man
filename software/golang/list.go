/*
 * Copyright © 2025 Piotr Kozak <piotrkrzysztofkozak@gmail.com>
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
package golang

import (
	"encoding/json"
	"fmt"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/console"
	"io"
	"net/http"
	"runtime"
	"strings"
)

type ApiPackagesPerVersion struct {
	Version string       `json:"version"`
	Files   []ApiPackage `json:"files"`
}

type ApiPackage struct {
	Version  string `json:"version"`
	Arch     string `json:"arch"`
	Os       string `json:"os"`
	Filename string `json:"filename"`
	Sha256   string `json:"sha256"`
	Kind     string `json:"kind"`
	Size     int64  `json:"size"`
}

type Package struct {
	Version      string
	Filename     string
	DownloadLink string
	Type         domain.Type
	Sha256       string
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
	var packagesPerVersion []ApiPackagesPerVersion
	err = json.NewDecoder(resp.Body).Decode(&packagesPerVersion)
	if err != nil {
		console.Fatal(err)
	}
	return supportedPackages(&packagesPerVersion, runtime.GOOS, runtime.GOARCH)
}

func supportedPackages(packagesPerVersions *[]ApiPackagesPerVersion, goOpSystem, goarch string) []Package {
	result := make([]Package, 0)
	for _, packagesPerVersion := range *packagesPerVersions {
		for _, apiPackage := range packagesPerVersion.Files {
			if apiPackage.Kind == "archive" && apiPackage.Arch == goarch && apiPackage.Os == goOpSystem {
				version, _ := strings.CutPrefix(packagesPerVersion.Version, "go")
				version = strings.Replace(version, "rc", ".rc", 1)
				version = strings.Replace(version, "beta", ".beta", 1)
				p := Package{
					Version:      version,
					Filename:     apiPackage.Filename,
					Type:         toType(apiPackage.Filename),
					DownloadLink: fmt.Sprintf("%s/%s", DistURL, apiPackage.Filename),
					Sha256:       apiPackage.Sha256,
				}
				result = append(result, p)
			}
		}
	}
	return result
}

func toType(filename string) domain.Type {
	if strings.HasSuffix(filename, ".zip") {
		return domain.ZIP
	} else if strings.HasSuffix(filename, ".tar.gz") || strings.HasSuffix(filename, ".tgz") {
		return domain.TAR_GZ
	} else {
		return domain.UNKNOWN
	}
}
