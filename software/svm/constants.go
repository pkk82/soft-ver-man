package svm

import "github.com/pkk82/soft-ver-man/util/github"

const RepoOwner = "pkk82"
const RepoName = "soft-ver-man"

const PageSize = 1000

const Name = "soft-ver-man"

const FileName = "svm"

func ReleasesURL() string {
	return github.URL(RepoOwner, RepoName)
}
