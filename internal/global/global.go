package global

import "fmt"

const GITHUB_LINK = "https://github.com/marcuswhybrow/ray-peat-rodeo"
const SPONSOR_LINK = "https://github.com/sponsors/marcuswhybrow"

func GitHubNewIssueLink() string {
	return GITHUB_LINK + "/issues/new"
}

// URL to open a GitHub issue
func GitHubIssueLink(id int) string {
	return GITHUB_LINK + "/issues/" + fmt.Sprint(id)
}

// URL to edit a file in the GitHub repository
func GitHubEditLink(filePath string) string {
	return GITHUB_LINK + "/edit/main/" + filePath
}
