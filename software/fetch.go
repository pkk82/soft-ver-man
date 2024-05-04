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

package software

import (
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/console"
	"github.com/pkk82/soft-ver-man/util/download"
	"path/filepath"
	"runtime"
)

func Fetch(plugin domain.Plugin, inputVersion, softwareDownloadDir string, verifyChecksum bool) (domain.FetchedPackage, error) {

	var version domain.Version
	var asset domain.Asset

	assets, err := plugin.GetAvailableAssets()
	if err == nil {
		versions := make([]string, len(assets))
		for i, v := range assets {
			versions[i] = v.Version
		}
		var index int
		version, index, err = domain.FindVersion(inputVersion, versions)
		if err != nil {
			return domain.FetchedPackage{}, err
		}
		asset = assets[index]
	} else {
		version, err = domain.NewVersion(inputVersion)
		if err != nil {
			return domain.FetchedPackage{}, err
		}
		downloadUrl, extension := plugin.CalculateDownloadUrl(version, runtime.GOOS, runtime.GOARCH)

		asset = domain.Asset{Url: downloadUrl, Type: extension, Version: inputVersion, Name: plugin.Name}
	}

	pluginDir := filepath.Join(softwareDownloadDir, plugin.Name)

	filename := plugin.CalculateDownloadedFileName(asset)

	fetchedPackagePath := download.FetchFile(asset.Url, pluginDir, filename)
	fetchedPackage := domain.FetchedPackage{Version: version, FilePath: fetchedPackagePath, Type: asset.Type}

	if verifyChecksum {
		err = plugin.VerifyChecksum(asset, fetchedPackage)
		if err != nil {
			return domain.FetchedPackage{}, err
		} else {
			console.Info("Checksum verified: OK")
		}
	}

	return fetchedPackage, nil

}
