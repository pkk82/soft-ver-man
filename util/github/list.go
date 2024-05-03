package github

import (
	"encoding/json"
	"fmt"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/console"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

func ListReleases(repoOwner, repoName string, pageSize int, assetNamePredicate func(string) bool) error {

	allAssets, err := GetSupportedAssets(repoOwner, repoName, pageSize, assetNamePredicate)
	if err != nil {
		return err
	}
	for _, asset := range allAssets {
		console.Info(asset.Version.Value)
	}
	return nil
}

func GetSupportedAssets(repoOwner, repoName string, pageSize int, predicate func(string) bool) ([]Asset, error) {
	nextPageUrl := URL(repoOwner, repoName) + fmt.Sprintf("?per_page=%d", pageSize)

	var allPackages []Asset

	var err error
	var packages []JsonRelease
	for {

		packages, nextPageUrl, err = getPageOfSupportedReleases(nextPageUrl)
		if err != nil {
			return nil, err
		}
		for _, pkg := range packages {
			for _, asset := range pkg.Assets {
				v, err := domain.NewVersion(pkg.Version)
				if err == nil && predicate(asset.Name) {
					compilerAsset := Asset{
						Version: v,
						Name:    asset.Name,
						Url:     asset.Url,
						Type:    toType(asset.ContentType, asset.Name),
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

func toType(contentType string, name string) domain.Type {
	switch contentType {
	case "application/zip":
		return domain.ZIP
	case "application/x-gzip":
		return domain.TAR_GZ
	case "raw":
		return domain.RAW
	}
	extension := filepath.Ext(name)
	switch extension {
	case ".zip":
		return domain.ZIP
	case ".tar.gz":
		return domain.TAR_GZ
	}
	return domain.UNKNOWN
}

type JsonRelease struct {
	Version string      `json:"tag_name"`
	Assets  []JsonAsset `json:"assets"`
}

func (p JsonRelease) version() (domain.Version, error) {
	return domain.NewVersion(p.Version)
}

type JsonAsset struct {
	Url         string `json:"browser_download_url"`
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
}

type Asset struct {
	Version domain.Version
	Name    string
	Url     string
	Type    domain.Type
}

func getPageOfSupportedReleases(url string) ([]JsonRelease, string, error) {

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
	var releases []JsonRelease
	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		return nil, "", err
	}

	linkHeader := resp.Header.Get("Link")

	return releases, nextPageLink(linkHeader), nil
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
