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

import "sort"

type VersionGranularity string

const (
	VersionGranularityMajor VersionGranularity = "MAJOR"
	VersionGranularityMinor VersionGranularity = "MINOR"
)

type ExtractStrategy string

// how to extract archive
const (

	// UseCompressedDirOrArchiveName - if archive contains only one top level directory, it will be used as target directory,
	// otherwise archive name without extension will be used
	UseCompressedDirOrArchiveName ExtractStrategy = "default"

	// ReplaceCompressedDirWithArchiveName - archive name without extension will be used as target directory,
	// and it will replace archive top level directory if there is only one
	ReplaceCompressedDirWithArchiveName ExtractStrategy = "archive_replace"
)

type Plugin struct {
	Name                        string
	ExecutableRelativePath      string
	VersionGranularity          VersionGranularity
	ExtractStrategy             ExtractStrategy
	RawExecutableName           string
	CalculateDownloadUrl        func(version Version, os, arch string) (string, Type)
	CalculateDownloadedFileName func(asset Asset) string
	PostInstall                 func(installedPackage InstalledPackage) error
	PostUninstall               func(version Version) error
	VerifyChecksum              func(asset Asset, fetchedPackage FetchedPackage) error
	GetAvailableAssets          func() ([]Asset, error)
}

var mainRegistry = make(map[string]Plugin)

func Register(plugin Plugin) {
	mainRegistry[plugin.Name] = plugin
}

func GetPlugin(name string) Plugin {
	return mainRegistry[name]
}

func GetPlugins() []Plugin {
	plugins := make([]Plugin, 0)
	for _, plugin := range mainRegistry {
		plugins = append(plugins, plugin)
	}
	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Name < plugins[j].Name
	})
	return plugins
}
