package software

import "github.com/pkk82/soft-ver-man/domain"

type Plugin struct {
	Name          string
	PostUninstall func(version domain.Version) error
}

var mainRegistry = make(map[string]Plugin)

func Register(plugin Plugin) {
	mainRegistry[plugin.Name] = plugin
}

func GetPlugin(name string) Plugin {
	return mainRegistry[name]
}
