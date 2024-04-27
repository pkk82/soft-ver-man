package kotlin

import "github.com/pkk82/soft-ver-man/util/github"

const RepoOwner = "JetBrains"
const RepoName = "kotlin"

const PageSize = 100

const Name = "kotlin"

var Aliases = []string{"kotlin", "kt"}

func ReleasesURL() string {
	return github.URL(RepoOwner, RepoName)
}
