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
package shell

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/config"
	"github.com/pkk82/soft-ver-man/history"
	"github.com/pkk82/soft-ver-man/pack"
	"github.com/pkk82/soft-ver-man/version"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type helper struct {
	path    string
	version version.Version
}

func initVariables(finder DirFinder, history history.PackageHistory) error {
	name := history.Name

	homeDir, err := finder.HomeDir()
	if err != nil {
		return err
	}

	softDir, err := finder.SoftDir()
	if err != nil {
		return err
	}

	exportLine := fmt.Sprintf("export SVM_DIR=%v", softDir)
	err = assertFileWithContent(path.Join(homeDir, config.HomeConfigDir, config.RcFile), exportLine,
		[]string{exportLine})
	if err != nil {
		return err
	}
	initLine := bashToLoad(makeRcName(name))
	err = assertFileWithContent(path.Join(homeDir, config.HomeConfigDir, config.RcFile), initLine,
		[]string{initLine})
	if err != nil {
		return err
	}

	installedPackages, err := convertHistoryToInstalledPackages(history)
	if err != nil {
		return err
	}

	installedPackagesPerMajorVersions := convertToMap(installedPackages)

	majorVersions := sortedMajorVersions(installedPackagesPerMajorVersions)
	lines := make([]string, 0)

	packageDirVarName := fmt.Sprintf("%v_DIR", strings.ToUpper(name))
	lines = append(lines, fmt.Sprintf("export %v=\"$SOFT_DIR/%v\"", packageDirVarName, name))

	for _, v := range majorVersions {
		latestMajor := installedPackagesPerMajorVersions[v][0]
		_, dirName := filepath.Split(latestMajor.Path)
		lines = append(lines, fmt.Sprintf("export %v_%v_HOME=\"$%v/%v\"", strings.ToUpper(name), v, packageDirVarName, dirName))
	}

	mainTime := int64(0)
	mainLine := ""
	for _, ip := range installedPackages {
		if ip.Main && ip.InstalledOn > mainTime {
			_, dirName := filepath.Split(ip.Path)
			mainLine = fmt.Sprintf("export %v_HOME=\"$%v/%v\"", strings.ToUpper(name), packageDirVarName, dirName)
			mainTime = ip.InstalledOn
		}
	}
	if mainLine == "" {
		_, dirName := filepath.Split(installedPackagesPerMajorVersions[majorVersions[len(majorVersions)-1]][0].Path)
		mainLine = fmt.Sprintf("export %v_HOME=\"$%v/%v\"", strings.ToUpper(name), packageDirVarName,
			dirName)
	}

	lines = append(lines, mainLine)

	lines = append(lines, fmt.Sprintf("export PATH=\"$%v_HOME/bin:$PATH\"", strings.ToUpper(name)))

	err = overrideFileWithContent(path.Join(homeDir, config.HomeConfigDir, makeRcName(name)), lines)
	if err != nil {
		return err
	}

	return nil

}

func convertHistoryToInstalledPackages(history history.PackageHistory) ([]pack.InstalledPackage, error) {
	installedPackages := make([]pack.InstalledPackage, 0)
	for _, item := range history.Items {
		v, err := version.NewVersion(item.Version)
		if err != nil {
			return nil, err
		}
		installedPackages = append(installedPackages, pack.InstalledPackage{
			Version:     v,
			Path:        item.Path,
			Main:        item.Main,
			InstalledOn: item.InstalledOn,
		})
	}
	return installedPackages, nil
}

func convertToMap(installedPackages []pack.InstalledPackage) map[int][]pack.InstalledPackage {
	packagesPerMajor := make(map[int][]pack.InstalledPackage)
	for _, ip := range installedPackages {
		packagesPerMajor[ip.Version.Major()] = append(packagesPerMajor[ip.Version.Major()], ip)
	}
	for _, v := range packagesPerMajor {
		sort.Slice(v, func(i, j int) bool {
			return version.CompareDesc(v[i].Version, v[j].Version)
		})
	}
	return packagesPerMajor
}

func sortedMajorVersions(perVersion map[int][]pack.InstalledPackage) []int {
	majors := make([]int, 0)
	for k := range perVersion {
		majors = append(majors, k)
	}
	sort.Slice(majors, func(i, j int) bool {
		return majors[i] < majors[j]
	})
	return majors
}
