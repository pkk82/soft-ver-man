package config

import (
	"github.com/pkk82/soft-ver-man/history"
	"github.com/spf13/viper"
)

func ReadHistoryConfig(name string) (history.PackageHistory, error) {
	if !viper.IsSet(name + PackageHistorySuffix) {
		return history.PackageHistory{Name: name, Items: []history.PackageHistoryItem{}}, nil
	}
	repr := viper.GetString(name + PackageHistorySuffix)
	return history.ParseHistory(repr)
}

func WriteHistoryConfig(history history.PackageHistory) error {
	json := history.Serialize()
	viper.Set(history.Name+PackageHistorySuffix, json)
	return viper.WriteConfig()
}
