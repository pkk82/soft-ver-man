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
