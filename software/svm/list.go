package svm

import (
	"github.com/pkk82/soft-ver-man/util/github"
	"runtime"
	"strings"
)

func List() error {
	return github.ListReleases(RepoOwner, RepoName, PageSize, predicate)
}

var predicate = func(name string) bool {
	return strings.Contains(name, runtime.GOARCH+"-"+runtime.GOOS)
}

func getSupportedPackages() ([]github.Asset, error) {
	return github.GetSupportedAssets(RepoOwner, RepoName, PageSize, predicate)
}
