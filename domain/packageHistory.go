package domain

import (
	"encoding/json"
)

type PackageHistory struct {
	Name  string               `json:"name"`
	Items []PackageHistoryItem `json:"items"`
}

type PackageHistoryItem struct {
	Version     string `json:"version"`
	Path        string `json:"path"`
	InstalledOn int64  `json:"installedOn"`
	Main        bool   `json:"main"`
}

func ParseHistory(repr string) (PackageHistory, error) {
	if repr == "" {
		return PackageHistory{}, nil
	}

	var history PackageHistory

	err := json.Unmarshal([]byte(repr), &history)
	if err != nil {
		return PackageHistory{}, err
	}

	return history, nil

}

func (ph *PackageHistory) Serialize() string {
	serialized, err := json.Marshal(ph)
	if err != nil {
		return ""
	}
	return string(serialized)
}

func (ph *PackageHistory) Add(installedPackage InstalledPackage) *PackageHistory {
	ph.Items = append(ph.Items, PackageHistoryItem{
		Version:     installedPackage.Version.Value,
		Path:        installedPackage.Path,
		InstalledOn: installedPackage.InstalledOn,
		Main:        false,
	})
	return ph
}

func (ph *PackageHistory) IsInstalled(version Version) bool {
	for _, item := range ph.Items {
		if item.Version == version.Value {
			return true
		}
	}
	return false
}
