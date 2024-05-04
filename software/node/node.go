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
package node

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/download"
	"github.com/pkk82/soft-ver-man/util/pgp"
	"path/filepath"
)

func init() {
	var plugin = domain.Plugin{
		Name:               Name,
		EnvNamePrefix:      EnvNamePrefix,
		GetAvailableAssets: getAvailableAssets,
		VerifyChecksum:     verifyChecksum,
		PostUninstall: func(version domain.Version) error {
			return nil
		},
		PostInstall: func(installedPackage domain.InstalledPackage) error {
			return nil
		},
		CalculateDownloadedFileName: calculateDownloadFileName,
		ExtractStrategy:             domain.UseCompressedDirOrArchiveName,
		ExecutableRelativePath:      "bin",
		VersionGranularity:          domain.VersionGranularityMajor,
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
		assets[i] = domain.Asset{
			Name:            p.FileName,
			Version:         p.Version,
			Url:             p.DownloadLink(),
			Type:            p.Type,
			ExtraProperties: map[string]string{"sumsLink": p.SumsLink(), "sumsSigLink": p.SumsSigLink()},
		}
	}
	return assets, nil
}

func verifyChecksum(asset domain.Asset, fetchedPackage domain.FetchedPackage) error {
	nodeDownloadDir := filepath.Dir(fetchedPackage.FilePath)
	softwareDownloadDir := filepath.Dir(nodeDownloadDir)
	version := fetchedPackage.Version
	publicKeyPaths := fetchPGPKeys(softwareDownloadDir)
	shaSumFileName := fmt.Sprintf("%v-%v-%v", Name, version.Value, ShaSumFileName)
	shaSumFilePath := download.FetchFileSilently(asset.ExtraProperties["sumsLink"], nodeDownloadDir, shaSumFileName)
	shaSumSigFileName := fmt.Sprintf("%v-%v-%v", Name, version.Value, ShaSumSigFileName)
	shaSumSigFilePath := download.FetchFileSilently(asset.ExtraProperties["sumsSigLink"], nodeDownloadDir, shaSumSigFileName)
	pgp.VerifySignature(shaSumFilePath, shaSumSigFilePath, publicKeyPaths)
	return verifySha(fetchedPackage.FilePath, shaSumFilePath)
}
