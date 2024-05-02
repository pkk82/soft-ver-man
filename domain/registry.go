package domain

type VersionGranularity string

const (
	VersionGranularityMajor VersionGranularity = "MAJOR"
	VersionGranularityMinor VersionGranularity = "MINOR"
)

type Plugin struct {
	Name                        string
	ExecutableRelativePath      string
	VersionGranularity          VersionGranularity
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
