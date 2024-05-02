package software

import (
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/download"
	"path/filepath"
	"runtime"
)

func fetch(plugin domain.Plugin, inputVersion, softwareDownloadDir string) (domain.FetchedPackage, error) {
	version, err := domain.NewVersion(inputVersion)
	if err != nil {
		return domain.FetchedPackage{}, err
	}
	pluginDir := filepath.Join(softwareDownloadDir, plugin.Name)
	downloadUrl, extension := plugin.CalculateDownloadUrl(version, runtime.GOOS, runtime.GOARCH)
	filename := plugin.CalculateDownloadedFileName(version, extension)
	fetchedPackagePath := download.FetchFile(downloadUrl, pluginDir, filename)
	return domain.FetchedPackage{Version: version, FilePath: fetchedPackagePath, Type: extension}, nil
}
