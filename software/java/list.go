/*
 * Copyright © 2024 Piotr Kozak <piotrkrzysztofkozak@gmail.com>
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

package java

import (
	"encoding/json"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/console"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

func getSupportedPackages() []Package {
	var allPackages []Package

	pageNo := 1
	for {
		pagination, packages := getPageOfSupportedPackages(pageNo)
		for _, pkg := range packages {
			allPackages = append(allPackages, pkg)
		}
		if pagination.NextPage == 0 {
			break
		}
		pageNo = pagination.NextPage
	}
	return allPackages
}

type Package struct {
	Name             string `json:"name"`
	Latest           bool   `json:"latest"`
	DownloadUrl      string `json:"download_url"`
	JavaVersion      []int  `json:"java_version"`
	AvailabilityType string `json:"availability_type"`
	Id               string `json:"package_uuid"`
}

func (p Package) version() string {
	return strings.Join(intSliceToStringSlice(p.JavaVersion), ".")
}

func (p Package) packagingType() domain.Type {
	if strings.HasSuffix(p.DownloadUrl, ".zip") {
		return domain.ZIP
	}
	return domain.TAR_GZ
}

type Pagination struct {
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
	FirstPage  int `json:"first_page"`
	LastPage   int `json:"last_page"`
	Page       int `json:"page"`
	NextPage   int `json:"next_page"`
}

func getPageOfSupportedPackages(pageNo int) (Pagination, []Package) {
	url := PackagesAPIURL + "?page=" + strconv.Itoa(pageNo) +
		"&page_size=" + strconv.Itoa(PageSize) +
		"&javafx_bundled=false" +
		"&crac_supported=false" +
		"&release_status=ga" +
		"&java_package_type=jdk" +
		"&os=" + toOs(runtime.GOOS) +
		"&arch=" + toArch(runtime.GOARCH) +
		"&archive_type=" + toType(runtime.GOOS)
	resp, err := http.Get(url)

	var pagination Pagination
	err = json.Unmarshal([]byte(resp.Header.Get("X-Pagination")), &pagination)
	if err != nil {
		console.Fatal(err)
	}

	if err != nil {
		console.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			console.Fatal(err)
		}
	}(resp.Body)
	var packages []Package
	err = json.NewDecoder(resp.Body).Decode(&packages)
	if err != nil {
		console.Fatal(err)
	}
	return pagination, packages
}

func toOs(goOpSystem string) string {
	if goOpSystem == "darwin" {
		return "macos"
	}
	if goOpSystem == "linux" {
		return "linux-glibc"
	}
	return goOpSystem
}

func toArch(goArch string) string {
	if goArch == "386" {
		return "x86"
	}
	return goArch
}

func toType(goOpSystem string) string {
	if goOpSystem == "linux" {
		return domain.TAR_GZ
	}
	return domain.ZIP
}

func intSliceToStringSlice(intSlice []int) []string {
	stringSlice := make([]string, len(intSlice))

	for i, num := range intSlice {
		stringSlice[i] = strconv.Itoa(num)
	}

	return stringSlice
}
