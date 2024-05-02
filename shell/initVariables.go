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
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/file"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

func initVariables(finder DirFinder, history domain.PackageHistory, executableDirName string, granularity domain.VersionGranularity) error {
	name := history.Name

	homeDir, err := finder.HomeDir()
	if err != nil {
		return err
	}

	err = initSvmRcFile(name, homeDir, finder)
	if err != nil {
		return err
	}

	installedPackages, err := history.ToInstalledPackages()
	if err != nil {
		return err
	}

	switch granularity {
	case domain.VersionGranularityMajor:
		err = initSpecificRcRileWithMajorVariables(installedPackages, name, homeDir, executableDirName)
		if err != nil {
			return err
		}
	case domain.VersionGranularityMinor:
		err = initSpecificRcRileWithMinorVariables(installedPackages, name, homeDir, executableDirName)
		if err != nil {
			return err
		}
	}

	return nil

}

func initSvmRcFile(name, homeDir string, finder DirFinder) error {

	softDir, err := finder.SoftDir()
	if err != nil {
		return err
	}

	exportLine := exportVariable(config.VarNameSvmSoftDir, softDir)
	err = assertFileWithContent(path.Join(homeDir, config.HomeConfigDir, config.RcFile), exportLine,
		[]string{exportLine})
	if err != nil {
		return err
	}
	initLine := bashToLoad(rcName(name))
	err = assertFileWithContent(path.Join(homeDir, config.HomeConfigDir, config.RcFile), initLine,
		[]string{initLine})
	if err != nil {
		return err
	}
	return nil
}

func initSpecificRcRileWithMajorVariables(installedPackages []domain.InstalledPackage, name, homeDir, executableDirName string) error {
	variablePartName := variablePartName(name)

	installedPackagesPerMajorVersions, err := toSortedMapPerVersion(installedPackages, MajorRounder)
	if err != nil {
		return err
	}

	lines := make([]string, 0)

	packageDirVarName := fmt.Sprintf(config.VarNameSvmSoftPackageDirTemplate, variablePartName)
	lines = append(lines, exportRefPathVariable(packageDirVarName, config.VarNameSvmSoftDir, name))

	for _, v := range installedPackagesPerMajorVersions.keys {
		get, _ := installedPackagesPerMajorVersions.get(v)
		latestMajor := get[0]
		_, dirName := filepath.Split(latestMajor.Path)
		lines = append(lines, exportHomeMajorVersionVariable(variablePartName, v, packageDirVarName, dirName))
	}

	mainTime := int64(0)
	mainLine := ""
	for _, ip := range installedPackages {
		if ip.Main && ip.InstalledOn > mainTime {
			_, dirName := filepath.Split(ip.Path)
			mainLine = exportHomeVariable(variablePartName, packageDirVarName, dirName)
			mainTime = ip.InstalledOn
		}
	}
	if mainLine == "" {
		last := installedPackagesPerMajorVersions.getLast()
		_, dirName := filepath.Split(last[0].Path)
		mainLine = exportHomeVariable(variablePartName, packageDirVarName, dirName)
	}

	lines = append(lines, mainLine)

	lines = append(lines, fmt.Sprintf("export PATH=\"$%v_HOME/%v:$PATH\"", variablePartName, executableDirName))

	err = file.OverrideFileWithContent(path.Join(homeDir, config.HomeConfigDir, rcName(name)), lines)
	if err != nil {
		return err
	}
	return nil
}

func initSpecificRcRileWithMinorVariables(installedPackages []domain.InstalledPackage, name, homeDir, executableDirName string) error {
	variablePartName := variablePartName(name)

	installedPackagesPerMajorVersions, err := toSortedMapPerVersion(installedPackages, MinorRounder)
	if err != nil {
		return err
	}

	lines := make([]string, 0)

	packageDirVarName := fmt.Sprintf(config.VarNameSvmSoftPackageDirTemplate, variablePartName)
	lines = append(lines, exportRefPathVariable(packageDirVarName, config.VarNameSvmSoftDir, name))

	for _, v := range installedPackagesPerMajorVersions.keys {
		minor, _ := installedPackagesPerMajorVersions.get(v)
		latestMinor := minor[0]
		_, dirName := filepath.Split(latestMinor.Path)
		lines = append(lines, exportHomeMinorVersionVariable(variablePartName, v, packageDirVarName, dirName))
	}

	mainTime := int64(0)
	mainLine := ""
	for _, ip := range installedPackages {
		if ip.Main && ip.InstalledOn > mainTime {
			_, dirName := filepath.Split(ip.Path)
			mainLine = exportHomeVariable(variablePartName, packageDirVarName, dirName)
			mainTime = ip.InstalledOn
		}
	}
	if mainLine == "" {
		last := installedPackagesPerMajorVersions.getLast()
		_, dirName := filepath.Split(last[0].Path)
		mainLine = exportHomeVariable(variablePartName, packageDirVarName, dirName)
	}

	lines = append(lines, mainLine)

	lines = append(lines, fmt.Sprintf("export PATH=\"$%v_HOME%v:$PATH\"", variablePartName, executableDirName))

	err = file.OverrideFileWithContent(path.Join(homeDir, config.HomeConfigDir, rcName(name)), lines)
	if err != nil {
		return err
	}
	return nil
}

type VersionRounder func(domain.InstalledPackage) (domain.Version, error)

func MinorRounder(ip domain.InstalledPackage) (domain.Version, error) {
	return domain.NewVersion(fmt.Sprintf("%d.%d", ip.Version.Major(), ip.Version.Minor()))
}

func MajorRounder(ip domain.InstalledPackage) (domain.Version, error) {
	return domain.NewVersion(fmt.Sprintf("%d", ip.Version.Major()))
}

func toSortedMapPerVersion(installedPackages []domain.InstalledPackage, versionRounder VersionRounder) (SortedMap, error) {
	packagesPerVersion := make(map[domain.Version][]domain.InstalledPackage)
	for _, ip := range installedPackages {
		roundedVersion, err := versionRounder(ip)
		if err != nil {
			return SortedMap{}, err
		}
		packagesPerVersion[roundedVersion] = append(packagesPerVersion[roundedVersion], ip)
	}
	sortedMap := newSortedMap()
	for v, p := range packagesPerVersion {
		sort.Slice(p, func(i, j int) bool {
			return domain.CompareDesc(p[i].Version, p[j].Version)
		})
		sortedMap.put(v, p)
	}
	return sortedMap, nil
}

type SortedMap struct {
	keys   []domain.Version
	values map[domain.Version][]domain.InstalledPackage
}

func newSortedMap() SortedMap {
	return SortedMap{
		keys:   make([]domain.Version, 0),
		values: make(map[domain.Version][]domain.InstalledPackage),
	}
}

func (sm *SortedMap) put(key domain.Version, value []domain.InstalledPackage) {
	sm.keys = append(sm.keys, key)
	sort.Slice(sm.keys, func(i, j int) bool {
		return domain.CompareAsc(sm.keys[i], sm.keys[j])
	})
	sm.values[key] = value
}

func (sm *SortedMap) getLast() []domain.InstalledPackage {
	lastKey := sm.keys[len(sm.keys)-1]
	return sm.values[lastKey]
}

func (sm *SortedMap) get(key domain.Version) ([]domain.InstalledPackage, bool) {
	value, ok := sm.values[key]
	return value, ok
}

func transformName(name string) string {
	if strings.Contains(name, "-") {
		parts := strings.Split(name, "-")
		short := ""
		for _, part := range parts {
			short += string(part[0])
		}
		return short
	} else {
		return name
	}
}

func rcName(name string) string {
	transformedName := transformName(name)
	return makeRcName(transformedName)
}

func variablePartName(name string) string {
	transformedName := transformName(name)
	return strings.ToUpper(transformedName)
}
