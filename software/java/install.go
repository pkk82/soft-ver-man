/*
Copyright © 2023 Piotr Kozak <piotrkrzysztofkozak@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package java

import (
	"errors"
	"github.com/pkk82/soft-ver-man/config"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/shell"
	"github.com/pkk82/soft-ver-man/util/archive"
	"github.com/spf13/viper"
	"path"
	"time"
)

func Install(fetchedPackage domain.FetchedPackage, softwareDir string) error {
	history, err := config.ReadHistoryConfig(Name)
	if err != nil {
		return err
	}

	if history.IsInstalled(fetchedPackage.Version) {
		return errors.New("Version " + fetchedPackage.Version.Value + " is already installed")
	}

	extractedPackage, err := archive.Extract(fetchedPackage, path.Join(softwareDir, Name), archive.TargetDirNameDefault)
	installedPackage := domain.InstalledPackage{
		Version:     extractedPackage.Version,
		Path:        extractedPackage.Path,
		Main:        false,
		InstalledOn: time.Now().UnixMilli(),
	}
	history.Add(installedPackage)
	if err != nil {
		return err
	}
	err = config.WriteHistoryConfig(history)
	if err != nil {
		return err
	}

	finder := shell.ProdDirFinder{SoftwareDir: viper.GetString(config.SoftwareDirKey)}
	err = shell.AddVariables(finder, history, "bin", domain.VersionGranularityMajor)
	if err != nil {
		return err
	}
	return nil
}
