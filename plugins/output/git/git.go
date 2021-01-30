package git

import (
	"errors"
	"fmt"
	"os"
	"strings"

	git "github.com/go-git/go-git/v5"
	config "github.com/go-git/go-git/v5/config"
	http "github.com/go-git/go-git/v5/plumbing/transport/http"
	logrus "github.com/sirupsen/logrus"

	common "github.com/parinithshekar/gitsink/common"
	plugins "github.com/parinithshekar/gitsink/plugins/interfaces"
	logger "github.com/parinithshekar/gitsink/wrap/logrus/v1"
)

var (
	log = logger.New()
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

	// Get authentication object for source
	sourceAccountID, sourceAccessToken, err := gitClient.input.Credentials()
	if err != nil {
		log.Errorf("Failed to fetch source credentials")
		return
	}
	sourceAuth := http.BasicAuth{
		Username: sourceAccountID,
		Password: sourceAccessToken,
	}

	for _, repo := range repos {

		var localRepo *git.Repository
		if _, err := os.Stat(repo.Slug); os.IsNotExist(err) {
			// Clone the repo
			co := git.CloneOptions{
				URL:  repo.Source,
				Auth: &sourceAuth,
			}
			co.Validate()
			localRepo, _ = git.PlainClone(repo.Slug, false, &co)
		} else {
			localRepo, err = git.PlainOpen(repo.Slug)
		}

		_, err = localRepo.CreateRemote(&config.RemoteConfig{
			Name: "target",
			URLs: []string{repo.Target},
		})
		if err != nil {
			log.WithFields(logrus.Fields{
				"integration": gitClient.integrationName,
				"repository":  repo.Slug,
				"error":       err.Error(),
			}).Errorf("Failed to set target remote")
		}

		failedTags, err := gitClient.SyncTags(repo, localRepo)
		if err != nil {
			if failedTags != nil {
				log.WithFields(logrus.Fields{
					"integration": gitClient.integrationName,
					"repository":  repo.Slug,
					"failedTags":  failedTags,
					"error":       err.Error(),
				}).Warningf("Some tags not synced")
			} else {
				log.WithFields(logrus.Fields{
					"integration": gitClient.integrationName,
					"repository":  repo.Slug,
					"error":       err.Error(),
				}).Warningf("Failed to sync tags")
			}
		}

		failedBranches, err := gitClient.SyncBranches(repo, localRepo)
		if err != nil {
			if failedBranches != nil {
				log.WithFields(logrus.Fields{
					"integration":    gitClient.integrationName,
					"repository":     repo.Slug,
					"failedBranches": failedBranches,
					"error":          err.Error(),
				}).Warningf("Some branches not synced")
			} else {
				log.WithFields(logrus.Fields{
					"integration": gitClient.integrationName,
					"repository":  repo.Slug,
					"error":       err.Error(),
				}).Warningf("Failed to sync branches")
			}
		}

		err = localRepo.DeleteRemote("target")
		if err != nil {
			log.WithFields(logrus.Fields{
				"integration": gitClient.integrationName,
				"repository":  repo.Slug,
				"error":       err.Error(),
			}).Warningf("Failed to remove old remote")
		}
	}

	// back to project root
	os.Chdir("../..")
}

// SyncTags individually syncs the tags from source remote to the target remote
func (gitClient Client) SyncTags(repo common.Repository, localRepo *git.Repository) ([]string, error) {

	var failedTags []string
	// Get authentication object for source
	sourceAccountID, sourceAccessToken, err := gitClient.input.Credentials()
	if err != nil {
		log.WithFields(logrus.Fields{
			"integration": gitClient.integrationName,
			"repository":  repo.Slug,
			"error":       err.Error(),
		}).Errorf("Failed to fetch source credentials")
		return nil, errors.New("Failed to sync tags")
	}
	sourceAuth := http.BasicAuth{
		Username: sourceAccountID,
		Password: sourceAccessToken,
	}

	// Get authentication object for target
	targetAccountID, targetAccessToken, err := gitClient.output.Credentials()
	if err != nil {
		log.WithFields(logrus.Fields{
			"integration": gitClient.integrationName,
			"repository":  repo.Slug,
			"error":       err.Error(),
		}).Errorf("Failed to fetch target credentials")
		return nil, errors.New("Failed to sync tags")
	}
	targetAuth := http.BasicAuth{
		Username: targetAccountID,
		Password: targetAccessToken,
	}

	// Fetch from origin
	fo := git.FetchOptions{
		RemoteName: "origin",
		Auth:       &sourceAuth,
	}
	fo.Validate()
	err = localRepo.Fetch(&fo)

	remotes, _ := localRepo.Remotes()
	for _, remote := range remotes {
		fmt.Println("\n", remote.String())
	}

	if err != nil && err.Error() != "already up-to-date" {
		log.WithFields(logrus.Fields{
			"integration": gitClient.integrationName,
			"repository":  repo.Slug,
			"error":       err.Error(),
		}).Warningf("Failed to get remote refs")
	}

	// Get list of origin tags
	origin, err := localRepo.Remote("origin")
	refs, err := origin.List(&git.ListOptions{
		Auth: &sourceAuth,
	})
	if err != nil {
		log.WithFields(logrus.Fields{
			"integration": gitClient.integrationName,
			"repository":  repo.Slug,
			"error":       err.Error(),
		}).Warningf("Failed to fetch tags from origin")
		return nil, errors.New("Failed to sync tags")
	}

	// Parse list of tag names
	var tags []string
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tags = append(tags, strings.SplitN(ref.Name().String(), "refs/tags/", 2)[1])
		}
	}

	// Sync tags
	for _, tag := range tags {

		// Build refspec
		tagRefspec := fmt.Sprintf("refs/tags/%v:refs/tags/%v", tag, tag)

		// Push tag to target remote
		po := git.PushOptions{
			RemoteName: "target",
			Auth:       &targetAuth,
			RefSpecs:   []config.RefSpec{config.RefSpec(tagRefspec)},
		}
		po.Validate()
		err = localRepo.Push(&po)

		// Report errors if any
		if (err != nil) && (err.Error() != "already up-to-date") {
			failedTags = append(failedTags, tag)
			log.WithFields(logrus.Fields{
				"integration": gitClient.integrationName,
				"repository":  repo.Slug,
				"tag":         tag,
				"error":       err.Error(),
			}).Errorf("Tag could not be synced")
		}

	}

	if len(failedTags) > 0 {
		return failedTags, errors.New("Some tags not synced")
	}
	return nil, nil
}

func reorderDefault(branches []string, keyBranch string) []string {
	if len(branches) == 0 || branches[0] == keyBranch {
		return branches
	}
	if branches[len(branches)-1] == keyBranch {
		branches = append([]string{keyBranch}, branches[:len(branches)-1]...)
		return branches
	}
	for i, branch := range branches {
		if branch == keyBranch {
			branches = append([]string{keyBranch}, append(branches[:i], branches[i+1:]...)...)
			break
		}
	}
	return branches
}

// SyncBranches individually syncs the branches from source remote to the target remote
func (gitClient Client) SyncBranches(repo common.Repository, localRepo *git.Repository) ([]string, error) {

	var failedBranches []string

	// Get authentication object for source
	sourceAccountID, sourceAccessToken, err := gitClient.input.Credentials()
	if err != nil {
		log.WithFields(logrus.Fields{
			"integration": gitClient.integrationName,
			"repository":  repo.Slug,
			"error":       err.Error(),
		}).Errorf("Failed to fetch source credentials")
		return nil, errors.New("Failed to sync tags")
	}
	sourceAuth := http.BasicAuth{
		Username: sourceAccountID,
		Password: sourceAccessToken,
	}

	// Get authentication object for target
	targetAccountID, targetAccessToken, err := gitClient.output.Credentials()
	if err != nil {
		log.WithFields(logrus.Fields{
			"integration": gitClient.integrationName,
			"repository":  repo.Slug,
			"error":       err.Error(),
		}).Errorf("Failed to fetch target credentials")
		return nil, errors.New("Failed to sync tags")
	}
	targetAuth := http.BasicAuth{
		Username: targetAccountID,
		Password: targetAccessToken,
	}

	// Fetch from origin
	fo := git.FetchOptions{
		RemoteName: "origin",
		Auth:       &sourceAuth,
	}
	fo.Validate()
	err = localRepo.Fetch(&fo)
	if err != nil && err.Error() != "already up-to-date" {
		log.WithFields(logrus.Fields{
			"integration": gitClient.integrationName,
			"repository":  repo.Slug,
			"error":       err.Error(),
		}).Warningf("Failed to fetch tags from origin")
	}

	// Get list of origin branches
	origin, err := localRepo.Remote("origin")
	refs, err := origin.List(&git.ListOptions{
		Auth: &sourceAuth,
	})
	if err != nil {
		log.WithFields(logrus.Fields{
			"integration": gitClient.integrationName,
			"repository":  repo.Slug,
			"error":       err.Error(),
		}).Warningf("Failed to get remote refs")
	}

	// Parse list of branch names
	var branches []string
	for _, ref := range refs {
		if ref.Name().IsBranch() {
			branches = append(branches, strings.SplitN(ref.Name().String(), "refs/heads/", 2)[1])
		}
	}

	// Reorder to push default branch from source first
	defaultBranch := "master"
	for _, ref := range refs {
		if ref.Strings()[0] == "HEAD" {
			defaultBranch = strings.SplitN(ref.Strings()[1], "ref: refs/heads/", 2)[1]
			break
		}
	}
	branches = reorderDefault(branches, defaultBranch)

	// Sync branches
	for _, branch := range branches {

		// Build refspec
		branchRefspec := fmt.Sprintf("refs/remotes/origin/%v:refs/heads/%v", branch, branch)

		// Push branch to target remote
		po := git.PushOptions{
			RemoteName: "target",
			Auth:       &targetAuth,
			RefSpecs:   []config.RefSpec{config.RefSpec(branchRefspec)},
		}
		po.Validate()
		err = localRepo.Push(&po)

		// Report errors if any
		if (err != nil) && (err.Error() != "already up-to-date") {
			failedBranches = append(failedBranches, branch)
			log.WithFields(logrus.Fields{
				"integration": gitClient.integrationName,
				"repository":  repo.Slug,
				"branch":      branch,
				"error":       err.Error(),
			}).Errorf("Branch could not be synced")
		}
	}

	if len(failedBranches) > 0 {
		return failedBranches, errors.New("Some branches not synced")
	}
	return nil, nil
}
