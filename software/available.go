package software

import (
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/console"
)

func Available(plugin domain.Plugin) {
	assets, err := plugin.GetAvailableAssets()
	if err == nil {
		for _, asset := range assets {
			console.Info(asset.Version)
		}
	} else {
		console.Error(err)
	}
}
