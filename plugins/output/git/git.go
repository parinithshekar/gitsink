package git

import (
	"os"

	"github.com/go-git/go-git/v5"

	"github.com/parinithshekar/github-migration-cli/common"
	plugins "github.com/parinithshekar/github-migration-cli/plugins/interfaces"
)

// Client struct has the output plugin associated with the integration
type Client struct {
	output *plugins.Output
}

// New returns a new git instance to perform git functions
func New(output *plugins.Output) (*Client) {
	gitClient := new(Client)
	gitClient.output = output
	return gitClient
}

// SyncRepos clones repositories locally and syncs 
func (gitClient Client) SyncRepos(repos []common.Repository) {
	if _, err := os.Stat("syncDirectory"); os.IsNotExist(err) {
		os.Mkdir("syncDirectory", os.ModeDir)
	}
	os.Chdir("syncDirectory")
	
	for _, repo := range repos {
		if _, err := os.Stat(repo.Slug); os.IsNotExist(err) {
			git.PlainClone(repo.Slug, false, &git.CloneOptions{ URL: repo.Source })
		}
		os.Chdir(repo.Slug)
		gitClient.SyncTags(repo)
		gitClient.SyncBranches(repo)
	}
}

// SyncTags individually syncs the tags from source remote to the target remote
func (gitClient Client) SyncTags(repo common.Repository) {
	// Sync Tags individually
}

// SyncBranches individually syncs the branches from source remote to the target remote
func (gitClient Client) SyncBranches(repo common.Repository) {
	// Sync Branches individually
}