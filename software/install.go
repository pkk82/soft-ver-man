package software

import (
	"errors"
	"github.com/pkk82/soft-ver-man/config"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/shell"
	"github.com/pkk82/soft-ver-man/util/archive"
	"github.com/pkk82/soft-ver-man/util/console"
	"github.com/pkk82/soft-ver-man/util/copy"
	"github.com/spf13/viper"
	"path"
	"time"
)

func Install(plugin domain.Plugin, inputVersion string, verifyChecksum bool) error {

	version, err := domain.NewVersion(inputVersion)
	if err != nil {
		return err
	}

	history, err := config.ReadHistoryConfig(plugin.Name)
	if err != nil {
		return err
	}

	if history.IsInstalled(version) {
		return errors.New("Version " + version.Value + " is already installed")
	}

	configuration, err := config.Get()
	if err != nil {
		return err
	}

	fetchedPackage, err := Fetch(plugin, inputVersion, configuration.SoftwareDownloadDir, verifyChecksum)
	if err != nil {
		console.Fatal(err)
	}

	var installedPackage domain.InstalledPackage
	if fetchedPackage.Type == domain.RAW {
		copiedPackage, err := copy.Copy(fetchedPackage, path.Join(configuration.SoftwareDir, plugin.Name, plugin.Name+"-"+fetchedPackage.Version.Value), plugin.RawExecutableName)
		if err != nil {
			return err
		}
		installedPackage = domain.InstalledPackage{
			Version:     copiedPackage.Version,
			Path:        copiedPackage.PathToFile,
			Main:        false,
			InstalledOn: time.Now().UnixMilli(),
		}
	} else {

		extractedPackage, err := archive.Extract(fetchedPackage, path.Join(configuration.SoftwareDir, plugin.Name), plugin.ExtractStrategy)
		if err != nil {
			return err
		}
		installedPackage = domain.InstalledPackage{
			Version:     extractedPackage.Version,
			Path:        extractedPackage.Path,
			Main:        false,
			InstalledOn: time.Now().UnixMilli(),
		}
	}
	history.Add(installedPackage)

	err = config.WriteHistoryConfig(history)
	if err != nil {
		return err
	}

	finder := shell.ProdDirFinder{SoftwareDir: viper.GetString(config.SoftwareDirKey)}
	err = shell.AddVariables(finder, history, plugin.ExecutableRelativePath, plugin.VersionGranularity)
	if err != nil {
		return err
	}

	err = plugin.PostInstall(installedPackage)
	if err != nil {
		return err
	}

	return nil
}
