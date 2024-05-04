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
