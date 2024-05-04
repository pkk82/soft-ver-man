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
	"errors"
	"github.com/pkk82/soft-ver-man/config"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/shell"
	"github.com/pkk82/soft-ver-man/util/archive"
	"github.com/pkk82/soft-ver-man/util/console"
	"github.com/pkk82/soft-ver-man/util/copy"
	"github.com/spf13/viper"
	"path"
	"time"
)

func Install(plugin domain.Plugin, inputVersion string, verifyChecksum bool) error {

	version, err := domain.NewVersion(inputVersion)
	if err != nil {
		return err
	}

	installedPackages, err := config.LoadInstalledPackages(plugin.Name)
	if err != nil {
		return err
	}

	if installedPackages.IsInstalled(version) {
		return errors.New("Version " + version.Value + " is already installed")
	}

	configuration, err := config.Get()
	if err != nil {
		return err
	}

	fetchedPackage, err := Fetch(plugin, inputVersion, configuration.SoftwareDownloadDir, verifyChecksum)
	if err != nil {
		console.Fatal(err)
	}

	var installedPackage domain.InstalledPackage
	if fetchedPackage.Type == domain.RAW {
		copiedPackage, err := copy.Copy(fetchedPackage, path.Join(configuration.SoftwareDir, plugin.Name, plugin.Name+"-"+fetchedPackage.Version.Value), plugin.RawExecutableName)
		if err != nil {
			return err
		}
		installedPackage = domain.InstalledPackage{
			Version:     copiedPackage.Version,
			Path:        copiedPackage.PathToFile,
			Main:        false,
			InstalledOn: time.Now().UnixMilli(),
		}
	} else {

		extractedPackage, err := archive.Extract(fetchedPackage, path.Join(configuration.SoftwareDir, plugin.Name), plugin.ExtractStrategy)
		if err != nil {
			return err
		}
		installedPackage = domain.InstalledPackage{
			Version:     extractedPackage.Version,
			Path:        extractedPackage.Path,
			Main:        false,
			InstalledOn: time.Now().UnixMilli(),
		}
	}
	installedPackages.Add(installedPackage)

	err = config.StoreInstalledPackages(installedPackages)
	if err != nil {
		return err
	}

	finder := shell.ProdDirFinder{SoftwareDir: viper.GetString(config.SoftwareDirKey)}
	err = shell.AddVariables(finder, installedPackages)
	if err != nil {
		return err
	}

	err = plugin.PostInstall(installedPackage)
	if err != nil {
		return err
	}

	return nil
}
