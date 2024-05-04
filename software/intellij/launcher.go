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
package intellij

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/config"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/shell"
	"github.com/pkk82/soft-ver-man/util/file"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path"
	"runtime"
	"strings"
)

func createLauncher(installedPackage domain.InstalledPackage) error {
	current, err := user.Current()
	if err != nil {
		return err
	}
	if runtime.GOOS == "linux" {
		version := installedPackage.Version.Value
		launcherFilename := "intellij-ultimate-" + version + ".desktop"
		where := path.Join(current.HomeDir, ".local", "share", "applications", launcherFilename)
		fmt.Printf("Creating launcher for Linux in %s\n", where)

		envVariables, err := prepareEnvVariables()
		if err != nil {
			return err
		}

		err = file.OverrideFileWithContent(where, []string{
			"[Desktop Entry]",
			"Name=IU " + version,
			"Comment=Intellij Ultimate " + version,
			"Exec=" + envVariables + " " + path.Join(installedPackage.Path, "bin", "idea.sh"),
			"Icon=" + path.Join(installedPackage.Path, "bin", "idea.png"),
			"Terminal=false",
			"Type=Application",
		})
		if err != nil {
			return err
		}

	}
	return nil
}

func prepareEnvVariables() (string, error) {
	allInstalledPackages, err := config.LoadAllInstalledPackages()
	if err != nil {
		return "", err
	}
	var lines []string
	for _, installedPackages := range allInstalledPackages {
		if len(installedPackages.Items) == 0 {
			continue
		}
		envVariables, err := installedPackages.PrepareEnvVariables(installedPackages.Plugin)
		if err != nil {
			return "", err
		}

		resolvedEnvVariables, err := envVariables.Resolve()
		if err != nil {
			return "", err
		}

		finder := shell.ProdDirFinder{SoftwareDir: viper.GetString(config.SoftwareDirKey)}
		softDir, err := finder.SoftDir()
		if err != nil {
			return "", err
		}
		lines = append(lines, shell.PrepareSvmSoftDirEnvVariable(softDir).ToEnv())
		for _, envVariable := range resolvedEnvVariables.Variables {
			lines = append(lines, envVariable.ToEnv())
		}
	}
	return strings.Join(lines, " "), nil
}

func deleteLauncher(version domain.Version) error {
	current, err := user.Current()
	if err != nil {
		return err
	}
	if runtime.GOOS == "linux" {
		versionValue := version.Value
		launcherFilename := "intellij-ultimate-" + versionValue + ".desktop"
		where := path.Join(current.HomeDir, ".local", "share", "applications", launcherFilename)
		fmt.Printf("Deleting launcher for Linux in %s\n", where)

		err := os.Remove(where)
		if err != nil {
			return err
		}
	}
	return nil
}
