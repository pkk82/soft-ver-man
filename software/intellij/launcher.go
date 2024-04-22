package intellij

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/file"
	"os/user"
	"path"
	"runtime"
)

func createLauncher(installedPackage domain.InstalledPackage) error {
	current, err := user.Current()
	if err != nil {
		return err
	}
	if runtime.GOOS == "linux" {
		version := installedPackage.Version.Value
		launcherFilename := "intellij-ultimate-" + version + ".desktop"
		where := path.Join(current.HomeDir, ".local", "share", "applications", launcherFilename)
		fmt.Printf("Creating launcher for Linux in %s\n", where)

		err := file.OverrideFileWithContent(where, []string{
			"[Desktop Entry]",
			"Name=IU " + version,
			"Comment=Intellij Ultimate " + version,
			// TODO provide shell variables
			"Exec=" + path.Join(installedPackage.Path, "bin", "idea.sh"),
			"Icon=" + path.Join(installedPackage.Path, "bin", "idea.png"),
			"Terminal=false",
			"Type=Application",
		})
		if err != nil {
			return err
		}

	}
	return nil
}
