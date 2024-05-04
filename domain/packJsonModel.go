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
	"encoding/json"
	"errors"
)

type installedPackagesJsonModel struct {
	Name  string                      `json:"name"`
	Items []installedPackageJsonModel `json:"items"`
}

type installedPackageJsonModel struct {
	Version     string `json:"version"`
	Path        string `json:"path"`
	InstalledOn int64  `json:"installedOn"`
	Main        bool   `json:"main"`
}

func (installedPackages *InstalledPackages) SerializeInstalledPackages() (string, error) {
	jsonModel := mapToJsonModel(installedPackages.Items)
	installedPackagesJsonModel := installedPackagesJsonModel{
		Name:  installedPackages.Plugin.Name,
		Items: jsonModel,
	}
	serialized, err := json.Marshal(installedPackagesJsonModel)
	if err != nil {
		return "", err
	}
	return string(serialized), nil
}

func DeserializeInstalledPackages(name, serialized string) (InstalledPackages, error) {
	if serialized == "" {
		return InstalledPackages{Plugin: GetPlugin(name), Items: make([]InstalledPackage, 0)}, nil
	}
	var model installedPackagesJsonModel
	err := json.Unmarshal([]byte(serialized), &model)
	if err != nil {
		return InstalledPackages{}, err
	}
	installedPackage, err := mapFromJsonModel(model.Items)
	if err != nil {
		return InstalledPackages{}, err
	}
	if model.Name != "" && model.Name != name {
		return InstalledPackages{}, errors.New("plugin name does not match")
	}
	return InstalledPackages{Plugin: GetPlugin(name), Items: installedPackage}, nil

}

func mapFromJsonModel(items []installedPackageJsonModel) ([]InstalledPackage, error) {
	result := make([]InstalledPackage, 0)
	for _, item := range items {
		version, err := NewVersion(item.Version)
		if err != nil {
			return nil, err
		}
		result = append(result, InstalledPackage{
			Version:     version,
			Path:        item.Path,
			InstalledOn: item.InstalledOn,
			Main:        item.Main,
		})
	}
	return result, nil
}

func mapToJsonModel(installedPackages []InstalledPackage) []installedPackageJsonModel {
	result := make([]installedPackageJsonModel, len(installedPackages))
	for index, item := range installedPackages {
		result[index] = installedPackageJsonModel{
			Version:     item.Version.Value,
			Path:        item.Path,
			InstalledOn: item.InstalledOn,
			Main:        item.Main,
		}
	}
	return result
}
