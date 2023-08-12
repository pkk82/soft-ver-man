package config

import (
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/spf13/viper"
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
