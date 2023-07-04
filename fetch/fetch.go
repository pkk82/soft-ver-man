package fetch

import "github.com/pkk82/soft-ver-man/version"

type FetchedPackage struct {
	Version  version.Version
	FilePath string
}
