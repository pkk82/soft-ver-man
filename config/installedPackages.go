package config

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/console"
	"github.com/spf13/viper"
	"os"
	"text/tabwriter"
)

func DisplayAllInstalledPackages() {

	plugins := domain.GetPlugins()
	for _, plugin := range plugins {
		fmt.Println()
		fmt.Println("Installed packages for: ", plugin.Name)
		DisplayInstalledPackages(plugin.Name)
	}

}

func DisplayInstalledPackages(name string) {

	installedPackages, err := LoadInstalledPackages(name)
	if err != nil {
		console.Fatal(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', tabwriter.Debug)
	if len(installedPackages.Items) == 0 {
		_, err = fmt.Fprintln(w, "There are no installed packages")
	} else {
		_, err = fmt.Fprintln(w, "Version\t Main\t Path")
	}
	if err != nil {
		console.Fatal(err)
	}
	for _, item := range installedPackages.Items {
		main := ""
		if item.Main {
			main = " Yes"
		}
		_, err = fmt.Fprintf(w, "%s\t %s\t %s\n", item.Version.Value, main, item.Path)
		if err != nil {
			console.Fatal(err)
		}
	}
	err = w.Flush()
	if err != nil {
		console.Fatal(err)
	}

}

func LoadInstalledPackages(name string) (domain.InstalledPackages, error) {
	json := readInstalledPackagesFromConfig(name)
	packages, err := domain.DeserializeInstalledPackages(name, json)
	if err != nil {
		return domain.InstalledPackages{}, err
	}

	return packages, nil
}

func StoreInstalledPackages(packages domain.InstalledPackages) error {

	installedPackages, err := packages.SerializeInstalledPackages()
	if err != nil {
		return err
	}

	err = writeInstalledPackagesToConfig(packages.Plugin.Name, installedPackages)
	if err != nil {
		return err
	}

	return nil

}

func readInstalledPackagesFromConfig(name string) string {
	if !viper.IsSet(name + InstalledPackagesSuffix) {
		return ""
	}
	return viper.GetString(name + InstalledPackagesSuffix)
}

func writeInstalledPackagesToConfig(name, json string) error {
	viper.Set(name+InstalledPackagesSuffix, json)
	return viper.WriteConfig()
}
