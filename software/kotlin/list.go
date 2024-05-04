package kotlin

import (
	"github.com/pkk82/soft-ver-man/util/github"
	"strings"
)

var predicate = func(name string) bool {
	return strings.HasPrefix(name, "kotlin-compiler") && strings.HasSuffix(name, ".zip")
}

func getSupportedPackages() ([]github.Asset, error) {
	return github.GetSupportedAssets(RepoOwner, RepoName, PageSize, predicate)
}
