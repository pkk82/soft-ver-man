package domain

type VersionGranularity string

const (
	VersionGranularityMajor VersionGranularity = "MAJOR"
	VersionGranularityMinor VersionGranularity = "MINOR"
)

type ExtractStrategy string

// how to extract archive
const (

	// UseCompressedDirOrArchiveName - if archive contains only one top level directory, it will be used as target directory,
	// otherwise archive name without extension will be used
	UseCompressedDirOrArchiveName ExtractStrategy = "default"

	// ReplaceCompressedDirWithArchiveName - archive name without extension will be used as target directory,
	// and it will replace archive top level directory if there is only one
	ReplaceCompressedDirWithArchiveName ExtractStrategy = "archive_replace"
)

type Plugin struct {
	Name                        string
	ExecutableRelativePath      string
	VersionGranularity          VersionGranularity
	ExtractStrategy             ExtractStrategy
	CalculateDownloadUrl        func(version Version, os, arch string) (string, Type)
	CalculateDownloadedFileName func(version Version, extension Type) string
	PostInstall                 func(installedPackage InstalledPackage) error
	PostUninstall               func(version Version) error
	VerifyChecksum              func(fetchedPackage FetchedPackage) error
}

var mainRegistry = make(map[string]Plugin)

func Register(plugin Plugin) {
	mainRegistry[plugin.Name] = plugin
}

func GetPlugin(name string) Plugin {
	return mainRegistry[name]
}
