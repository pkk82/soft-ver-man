package config

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/console"
	"github.com/spf13/viper"
	"os"
	"text/tabwriter"
)

func ReadHistoryConfig(name string) (domain.PackageHistory, error) {
	if !viper.IsSet(name + PackageHistorySuffix) {
		return domain.PackageHistory{Name: name, Items: []domain.PackageHistoryItem{}}, nil
	}
	repr := viper.GetString(name + PackageHistorySuffix)
	return domain.ParseHistory(repr)
}

func WriteHistoryConfig(history domain.PackageHistory) error {
	json := history.Serialize()
	viper.Set(history.Name+PackageHistorySuffix, json)
	return viper.WriteConfig()
}

func DisplayHistory(name string) {

	history, err := ReadHistoryConfig(name)
	if err != nil {
		console.Fatal(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', tabwriter.Debug)
	_, err = fmt.Fprintln(w, "Version\t Main\t Path")
	if err != nil {
		console.Fatal(err)
	}
	for _, item := range history.Items {
		main := ""
		if item.Main {
			main = " Yes"
		}
		_, err = fmt.Fprintf(w, "%s\t %s\t %s\n", item.Version, main, item.Path)
		if err != nil {
			console.Fatal(err)
		}
	}
	err = w.Flush()
	if err != nil {
		console.Fatal(err)
	}

}
