package domain

type EnvVariableGranularity string

const (
	EnvVariableGranularityMajor EnvVariableGranularity = "MAJOR"
	EnvVariableGranularityMinor EnvVariableGranularity = "MINOR"
)

type Plugin struct {
	Name                        string
	ExecutableRelativePath      string
	EnvVariableGranularity      EnvVariableGranularity
	CalculateDownloadUrl        func(version Version, os, arch string) (string, Type)
	CalculateDownloadedFileName func(version Version, extension Type) string
	PostInstall                 func(installedPackage InstalledPackage) error
	PostUninstall               func(version Version) error
}

var mainRegistry = make(map[string]Plugin)

func Register(plugin Plugin) {
	mainRegistry[plugin.Name] = plugin
}

func GetPlugin(name string) Plugin {
	return mainRegistry[name]
}
