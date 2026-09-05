/*
 * Copyright © 2026 Piotr Kozak <piotrkrzysztofkozak@gmail.com>
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
package gradle

import (
	"encoding/json"
	"fmt"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/console"
	"io"
	"net/http"
	"sort"
)

type ApiVersion struct {
	Version        string `json:"version"`
	Snapshot       bool   `json:"snapshot"`
	Nightly        bool   `json:"nightly"`
	ReleaseNightly bool   `json:"releaseNightly"`
	RcFor          string `json:"rcFor"`
	MilestoneFor   string `json:"milestoneFor"`
	Broken         bool   `json:"broken"`
	DownloadUrl    string `json:"downloadUrl"`
	Checksum       string `json:"checksum"`
}

type Package struct {
	Version      domain.Version
	DownloadLink string
	Sha256       string
}

func getSupportedPackages() ([]Package, error) {
	resp, err := http.Get(VersionsURL)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			console.Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status code for fetching versions %v", resp.StatusCode)
	}
	var apiVersions []ApiVersion
	err = json.NewDecoder(resp.Body).Decode(&apiVersions)
	if err != nil {
		return nil, err
	}
	return supportedPackages(apiVersions), nil
}

// supportedPackages keeps final releases only, newest first
func supportedPackages(apiVersions []ApiVersion) []Package {
	result := make([]Package, 0)
	for _, apiVersion := range apiVersions {
		if !isFinalRelease(apiVersion) {
			continue
		}
		version, err := domain.NewVersion(apiVersion.Version)
		if err != nil {
			continue
		}
		result = append(result, Package{
			Version:      version,
			DownloadLink: apiVersion.DownloadUrl,
			Sha256:       apiVersion.Checksum,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return domain.CompareDesc(result[i].Version, result[j].Version)
	})
	return result
}

// isFinalRelease excludes snapshots, nightlies, release candidates, milestones and broken builds;
// the version must also be plain numbers and dots, because the API marks early
// milestones such as 1.0-milestone-9 as regular releases
func isFinalRelease(apiVersion ApiVersion) bool {
	return !apiVersion.Snapshot &&
		!apiVersion.Nightly &&
		!apiVersion.ReleaseNightly &&
		!apiVersion.Broken &&
		apiVersion.RcFor == "" &&
		apiVersion.MilestoneFor == "" &&
		domain.ValidateVersion(apiVersion.Version) == nil
}
