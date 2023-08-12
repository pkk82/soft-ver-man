package kotlin

import (
	"encoding/json"
	"fmt"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/console"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func List() error {
	allPackages, err := getSupportedPackages()
	if err != nil {
		return err
	}
	for _, p := range allPackages {
		console.Info(p.Version.Value)
	}
	return nil
}

func getSupportedPackages() ([]CompilerAsset, error) {
	var allPackages []CompilerAsset

	nextPageUrl := fmt.Sprintf("%v?per_page=%d", ReleasesURL, PageSize)
	var err error
	var packages []Package
	for {

		packages, nextPageUrl, err = getPageOfSupportedPackages(nextPageUrl)
		if err != nil {
			return nil, err
		}
		for _, pkg := range packages {
			for _, asset := range pkg.Assets {
				v, err := domain.NewVersion(pkg.Version)
				if err == nil && strings.HasPrefix(asset.Name, "kotlin-compiler") && strings.HasSuffix(asset.Name, ".zip") {
					compilerAsset := CompilerAsset{
						Version: v,
						Name:    asset.Name,
						Url:     asset.Url,
						Type:    domain.ZIP,
					}
					allPackages = append(allPackages, compilerAsset)
				}
			}
		}
		if nextPageUrl == "" {
			break
		}
	}
	return allPackages, nil
}

type Package struct {
	Version string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Url  string `json:"browser_download_url"`
	Name string `json:"name"`
}

type CompilerAsset struct {
	Version domain.Version
	Name    string
	Url     string
	Type    domain.Type
}

func getPageOfSupportedPackages(url string) ([]Package, string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("HTTP status code for fetching releases %v", resp.StatusCode)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			console.Error(err)
		}
	}(resp.Body)
	var packages []Package
	err = json.NewDecoder(resp.Body).Decode(&packages)
	if err != nil {
		return nil, "", err
	}

	linkHeader := resp.Header.Get("Link")

	return packages, nextPageLink(linkHeader), nil
}

func nextPageLink(linkHeader string) string {
	linkHeaderPattern := regexp.MustCompile(`<(\S*)>; rel="next"`)
	pagesRemaining := linkHeader != "" && strings.Contains(linkHeader, `rel="next"`)
	if pagesRemaining {
		match := linkHeaderPattern.FindStringSubmatch(linkHeader)
		if len(match) > 1 {
			return match[1]
		}
	}
	return ""
}
