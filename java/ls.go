package java

import (
	"encoding/json"
	"github.com/pkk82/soft-ver-man/console"
	"github.com/pkk82/soft-ver-man/pack"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

func List() {
	var allPackages []Package

	pageNo := 1
	for {
		pagination, packages := getSupportedPackages(pageNo)
		for _, pkg := range packages {
			allPackages = append(allPackages, pkg)
		}
		if pagination.NextPage == 0 {
			break
		}
		pageNo = pagination.NextPage
	}
	for _, p := range allPackages {
		console.Info(strings.Join(intSliceToStringSlice(p.JavaVersion), "."))
	}
}

type Package struct {
	Name             string `json:"name"`
	Latest           bool   `json:"latest"`
	DownloadUrl      string `json:"download_url"`
	JavaVersion      []int  `json:"java_version"`
	AvailabilityType string `json:"availability_type"`
}

type Pagination struct {
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
	FirstPage  int `json:"first_page"`
	LastPage   int `json:"last_page"`
	Page       int `json:"page"`
	NextPage   int `json:"next_page"`
}

func getSupportedPackages(pageNo int) (Pagination, []Package) {
	url := PackagesAPIURL + "?page=" + strconv.Itoa(pageNo) +
		"&page_size=" + strconv.Itoa(PageSize) +
		"&javafx_bundled=false" +
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
		return pack.TAR_GZ
	}
	return pack.ZIP
}

func intSliceToStringSlice(intSlice []int) []string {
	stringSlice := make([]string, len(intSlice))

	for i, num := range intSlice {
		stringSlice[i] = strconv.Itoa(num)
	}

	return stringSlice
}
