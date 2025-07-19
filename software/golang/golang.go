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
package golang

import (
	"github.com/pkk82/soft-ver-man/domain"
	"strings"
)

func init() {
	var plugin = domain.Plugin{
		Name:               Name,
		EnvNamePrefix:      EnvNamePrefix,
		EnvNameSuffix:      EnvNameSuffix,
		GetAvailableAssets: getAvailableAssets,
		VerifyChecksum:     verifyChecksum,
		PostUninstall: func(version domain.Version) error {
			return nil
		},
		PostInstall: func(installedPackage domain.InstalledPackage) error {
			return nil
		},
		CalculateDownloadedFileName: calculateDownloadFileName,
		ExtractStrategy:             domain.ReplaceCompressedDirWithArchiveName,
		ExecutableRelativePath:      "bin",
		VersionGranularity:          domain.VersionGranularityMinor,
	}
	domain.Register(plugin)
}

func calculateDownloadFileName(asset domain.Asset) string {
	return asset.Name
}

func getAvailableAssets() ([]domain.Asset, error) {
	packages := getSupportedPackages()
	assets := make([]domain.Asset, len(packages))
	for i, p := range packages {
		version, _ := strings.CutPrefix(p.Version, "go")
		assets[i] = domain.Asset{
			Version:         version,
			Name:            p.Filename,
			Url:             p.DownloadLink,
			Type:            p.Type,
			ExtraProperties: map[string]string{"sha256": p.Sha256},
		}
	}
	return assets, nil
}

func verifyChecksum(asset domain.Asset, fetchedPackage domain.FetchedPackage) error {
	return verifySha(fetchedPackage.FilePath, asset.ExtraProperties["sha256"])
}
