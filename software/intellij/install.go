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
package intellij

import (
	"errors"
	"github.com/pkk82/soft-ver-man/config"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/shell"
	"github.com/pkk82/soft-ver-man/util/archive"
	"github.com/pkk82/soft-ver-man/util/console"
	"github.com/spf13/viper"
	"path"
	"time"
)

/*
	if err != nil {
		console.Fatal(err)
	}

installedPackage, err := intellij.Install(cmdUtil.config.SoftwareDir)

	if err != nil {
		console.Fatal(err)
	}

err = intellij.CreateLauncher(installedPackage)
*/
func Install(inputVersion string) error {

	version, err := domain.NewVersion(inputVersion)
	if err != nil {
		return err
	}

	history, err := config.ReadHistoryConfig(Name)
	if err != nil {
		return err
	}

	if history.IsInstalled(version) {
		return errors.New("Version " + version.Value + " is already installed")
	}

	configuration, err := config.Get()
	if err != nil {
		return err
	}

	fetchedPackage, err := fetch(inputVersion, configuration.SoftwareDownloadDir)
	if err != nil {
		console.Fatal(err)
	}

	extractedPackage, err := archive.Extract(fetchedPackage, path.Join(configuration.SoftwareDir, Name), archive.TargetDirNameArchiveReplace)
	if err != nil {
		return err
	}
	installedPackage := domain.InstalledPackage{
		Version:     extractedPackage.Version,
		Path:        extractedPackage.Path,
		Main:        false,
		InstalledOn: time.Now().UnixMilli(),
	}
	history.Add(installedPackage)

	err = config.WriteHistoryConfig(history)
	if err != nil {
		return err
	}

	finder := shell.ProdDirFinder{SoftwareDir: viper.GetString(config.SoftwareDirKey)}
	err = shell.AddVariables(finder, history, "/bin", shell.VariableGranularityMajor)
	if err != nil {
		return err
	}

	err = createLauncher(installedPackage)
	if err != nil {
		return err
	}

	return nil
}
