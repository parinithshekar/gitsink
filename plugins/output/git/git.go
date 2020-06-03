package git

import (
	"os"
	"fmt"
	"strings"

	git "github.com/go-git/go-git/v5"

	common "github.com/parinithshekar/gitsink/common"
	plugins "github.com/parinithshekar/gitsink/plugins/interfaces"
)

// Client struct has the output plugin associated with the integration
type Client struct {
	input           plugins.Input
	output          plugins.Output
	integrationName string
}

// New returns a new git instance to perform git functions
func New(input plugins.Input, output plugins.Output, integrationName string) *Client {
	gitClient := new(Client)

	gitClient.input = input
	gitClient.output = output

	integrationNameSplit := strings.Split(integrationName, " ")
	gitClient.integrationName = strings.Join(integrationNameSplit, "-")

	return gitClient
}

// SyncRepos clones repositories locally and syncs
func (gitClient Client) SyncRepos(repos []common.Repository) {

	// Make and enter syncDirectory if it does not exist
	if _, err := os.Stat("syncDirectory"); os.IsNotExist(err) {
		os.Mkdir("syncDirectory", 0777)
	}
	os.Chdir("syncDirectory")

	// Make and enter directory for the current integration
	if _, err := os.Stat(gitClient.integrationName); os.IsNotExist(err) {
		os.Mkdir(gitClient.integrationName, 0777)
	}
	os.Chdir(gitClient.integrationName)

	for _, repo := range repos {

		var localRepo *git.Repository
		if _, err := os.Stat(repo.Slug); os.IsNotExist(err) {

			// Get authenticated link for cloning repo
			sourceAccountID, sourceAccessToken := gitClient.input.Credentials()
			sourceLinkDomain := strings.SplitN(repo.Source, "//", 2)[1]
			authSourceLink := fmt.Sprintf("https://%v:%v@%v", sourceAccountID, sourceAccessToken, sourceLinkDomain)

			// Clone the repo
			localRepo, _ = git.PlainClone(repo.Slug, false, &git.CloneOptions{URL: authSourceLink})
		}
		os.Chdir(repo.Slug)

		// TODO: yet to handle and report errors in logs
		gitClient.SyncTags(repo, localRepo)
		gitClient.SyncBranches(repo, localRepo)

		// back to /syncDirectory/integration-name
		os.Chdir("..")
	}

	// back to project root
	os.Chdir("../..")
}

// SyncTags individually syncs the tags from source remote to the target remote
func (gitClient Client) SyncTags(repo common.Repository, localRepo *git.Repository) {
	// Sync Tags individually

	// Get authenticated link for fetching
	sourceAccountID, sourceAccessToken := gitClient.input.Credentials()
	sourceLinkDomain := strings.SplitN(repo.Source, "//", 2)[1]
	authSourceLink := fmt.Sprintf("https://%v:%v@%v", sourceAccountID, sourceAccessToken, sourceLinkDomain)
	authSourceLink = string(authSourceLink)

	// Get authenticated link for pushing
	targetAccountID, targetAccessToken := gitClient.output.Credentials()
	targetLinkDomain := strings.SplitN(repo.Target, "//", 2)[1]
	authTargetLink := fmt.Sprintf("https://%v:%v@%v", targetAccountID, targetAccessToken, targetLinkDomain)
	authTargetLink = string(authTargetLink)

	localRepo.Fetch(&git.FetchOptions{})

}

// SyncBranches individually syncs the branches from source remote to the target remote
func (gitClient Client) SyncBranches(repo common.Repository, localRepo *git.Repository) {
	// Sync Branches individually

	// Get authenticated link for fetching
	sourceAccountID, sourceAccessToken := gitClient.input.Credentials()
	sourceLinkDomain := strings.SplitN(repo.Source, "//", 2)[1]
	authSourceLink := fmt.Sprintf("https://%v:%v@%v", sourceAccountID, sourceAccessToken, sourceLinkDomain)
	authSourceLink = string(authSourceLink)

	// Get authenticated link for pushing
	targetAccountID, targetAccessToken := gitClient.output.Credentials()
	targetLinkDomain := strings.SplitN(repo.Target, "//", 2)[1]
	authTargetLink := fmt.Sprintf("https://%v:%v@%v", targetAccountID, targetAccessToken, targetLinkDomain)
	authTargetLink = string(authTargetLink)

	localRepo.Fetch(&git.FetchOptions{})
}
