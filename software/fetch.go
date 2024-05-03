package software

import (
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/console"
	"github.com/pkk82/soft-ver-man/util/download"
	"path/filepath"
	"runtime"
)

func Fetch(plugin domain.Plugin, inputVersion, softwareDownloadDir string, verifyChecksum bool) (domain.FetchedPackage, error) {

	var version domain.Version
	var asset domain.Asset

	assets, err := plugin.GetAvailableAssets()
	if err == nil {
		versions := make([]string, len(assets))
		for i, v := range assets {
			versions[i] = v.Version
		}
		var index int
		version, index, err = domain.FindVersion(inputVersion, versions)
		if err != nil {
			return domain.FetchedPackage{}, err
		}
		asset = assets[index]
	} else {
		version, err = domain.NewVersion(inputVersion)
		if err != nil {
			return domain.FetchedPackage{}, err
		}
		downloadUrl, extension := plugin.CalculateDownloadUrl(version, runtime.GOOS, runtime.GOARCH)

		asset = domain.Asset{Url: downloadUrl, Type: extension, Version: inputVersion, Name: plugin.Name}
	}

	pluginDir := filepath.Join(softwareDownloadDir, plugin.Name)

	filename := plugin.CalculateDownloadedFileName(asset)

	fetchedPackagePath := download.FetchFile(asset.Url, pluginDir, filename)
	fetchedPackage := domain.FetchedPackage{Version: version, FilePath: fetchedPackagePath, Type: asset.Type}

	if verifyChecksum {
		err = plugin.VerifyChecksum(asset, fetchedPackage)
		if err != nil {
			return domain.FetchedPackage{}, err
		} else {
			console.Info("Checksum verified: OK")
		}
	}

	return fetchedPackage, nil

}
