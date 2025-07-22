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
package domain

import (
	"errors"
	"fmt"
	"github.com/pkk82/soft-ver-man/util/collections"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type InstalledPackages struct {
	Plugin Plugin
	Items  []InstalledPackage
}

func (installedPackages *InstalledPackages) IsInstalled(version Version) bool {
	for _, item := range installedPackages.Items {
		if item.Version == version {
			return true
		}
	}
	return false
}

func (installedPackages *InstalledPackages) Add(installedPackage InstalledPackage) {
	if installedPackage.Main {
		for i, item := range installedPackages.Items {
			if item.Main {
				installedPackages.Items[i].Main = false
			}
		}
	}
	installedPackages.Items = append(installedPackages.Items, installedPackage)
}

func (installedPackages *InstalledPackages) RemoveByVersion(version Version) *InstalledPackage {
	newItems := make([]InstalledPackage, 0)
	var itemToRemove *InstalledPackage = nil
	for _, item := range installedPackages.Items {
		if item.Version != version {
			newItems = append(newItems, item)
		} else {
			copyItem := item
			itemToRemove = &copyItem
		}
	}
	installedPackages.Items = newItems
	return itemToRemove
}

func (installedPackages *InstalledPackages) FoundMain() (*InstalledPackage, error) {

	var items = make([]InstalledPackage, len(installedPackages.Items))
	copy(items, installedPackages.Items)

	if len(items) == 0 {
		return nil, errors.New("no packages found")
	}

	// find latest main
	sort.Slice(items, func(i, j int) bool {
		return items[i].InstalledOn > items[j].InstalledOn
	})

	for _, p := range items {
		if p.Main {
			return &p, nil
		}
	}

	// find latest version
	sort.Slice(items, func(i, j int) bool {
		return CompareDesc(items[i].Version, items[j].Version)
	})

	return &items[0], nil

}

func (installedPackages *InstalledPackages) PrepareEnvVariables(plugin Plugin) (EnvVariables, error) {

	var items = make([]InstalledPackage, len(installedPackages.Items))
	copy(items, installedPackages.Items)

	classifier := func(ip InstalledPackage) Version {
		version, err := ip.RoundVersion(plugin.VersionGranularity)
		if err != nil {
			panic(err)
		}
		return version
	}

	packagesByRoundedVersion := collections.GroupByAndCollect(items, classifier, transformer)
	roundedVersions := collections.Keys(packagesByRoundedVersion)
	sort.Slice(roundedVersions, func(i, j int) bool {
		return CompareAsc(roundedVersions[i], roundedVersions[j])
	})

	svmSoftDirVariable := EnvVariable{
		Name: VarNameSvmSoftDir,
	}

	refVariable := EnvVariable{
		Name:           fmt.Sprintf(VarNameSvmSoftPackageDirTemplate, plugin.EnvNamePrefix),
		PrefixVariable: &svmSoftDirVariable,
		SuffixValue:    plugin.Name,
	}

	envVariables := make([]EnvVariable, 0)
	envVariables = append(envVariables, refVariable)

	mainPackage, err := installedPackages.FoundMain()
	if err != nil {
		return EnvVariables{}, err
	}

	var mainVariableReference EnvVariable
	for _, roundedVersion := range roundedVersions {
		packages, _ := packagesByRoundedVersion[roundedVersion]
		firstPackage, found := packages.Get(0)
		if !found {
			panic("no package found")
		}

		_, dirName := filepath.Split(firstPackage.Path)
		envVariable := EnvVariable{
			PrefixVariable: &refVariable,
			SuffixValue:    dirName,
			Name:           versionedHomeVariable(plugin, roundedVersion),
		}

		for _, p := range packages.Values {
			if mainPackage.Version == p.Version {
				mainVariableReference = envVariable
			}
		}
		envVariables = append(envVariables, envVariable)
	}

	mainVariable := EnvVariable{
		Name:           homeVariable(plugin),
		PrefixVariable: &mainVariableReference,
	}
	envVariables = append(envVariables, mainVariable)

	return EnvVariables{
		Variables:              envVariables,
		MainVariable:           &mainVariable,
		ExecutableRelativePath: plugin.ExecutableRelativePath,
	}, nil

}

func transformer(ips []InstalledPackage) collections.SortedList[InstalledPackage] {
	return collections.NewSortedList(ips, func(i, j InstalledPackage) bool {
		return CompareDesc(i.Version, j.Version)
	})
}

func versionedHomeVariable(plugin Plugin, version Version) string {
	var v string
	switch plugin.VersionGranularity {
	case VersionGranularityMajor:
		v = strconv.Itoa(version.Major())
	case VersionGranularityMinor:
		v = fmt.Sprintf("%v_%v", version.Major(), version.Minor())

	}
	envNameSuffix := plugin.EnvNameSuffix
	if !strings.HasPrefix(envNameSuffix, "_") {
		envNameSuffix = "_" + envNameSuffix
	}
	return fmt.Sprintf("%v_%v%v", plugin.EnvNamePrefix, v, envNameSuffix)
}

func homeVariable(plugin Plugin) string {
	return fmt.Sprintf("%v%v", plugin.EnvNamePrefix, plugin.EnvNameSuffix)
}
