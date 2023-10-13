package github

import "fmt"

func URL(repoOwner, repoName string) string {
	return fmt.Sprintf(releasesURLTemplate, repoOwner, repoName)
}
