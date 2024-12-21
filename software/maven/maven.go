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
package maven

import (
	"errors"
	"fmt"
	"github.com/pkk82/soft-ver-man/domain"
	"strconv"
)

func init() {
	var plugin = domain.Plugin{
		Name:          Name,
		EnvNamePrefix: EnvNamePrefix,
		PostInstall: func(installedPackage domain.InstalledPackage) error {
			return nil
		},
		PostUninstall: func(version domain.Version) error {
			return nil
		},
		VerifyChecksum: func(asset domain.Asset, fetchedPackage domain.FetchedPackage) error {
			return errors.New("verify checksum not supported")
		},
		GetAvailableAssets: func() ([]domain.Asset, error) {
			return []domain.Asset{}, errors.New("get supported packages not supported")
		},
		CalculateDownloadUrl:        calculateDownloadUrl,
		CalculateDownloadedFileName: calculateDownloadedFileName,
		ExecutableRelativePath:      "bin",
		VersionGranularity:          domain.VersionGranularityMajor,
		ExtractStrategy:             domain.UseCompressedDirOrArchiveName,
	}
	domain.Register(plugin)
}

func calculateDownloadUrl(version domain.Version, _ string, _ string) (string, domain.Type) {
	extension := toExtension()
	url := fmt.Sprintf(DownloadURL, strconv.Itoa(version.Major()), version.Value, version.Value, extension)
	return url, extension
}

func toExtension() domain.Type {
	return domain.ZIP
}

func calculateDownloadedFileName(asset domain.Asset) string {
	return fmt.Sprintf("apache-maven-%s.%s", asset.Version, asset.Type)
}