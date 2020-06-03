package public

import (
	"os"
	"fmt"
	"strings"
	"context"
	"encoding/json"

	oauth2 "golang.org/x/oauth2"
	gjson "github.com/tidwall/gjson"
	github "github.com/google/go-github/v31/github"

	common "github.com/parinithshekar/gitsink/common"
	config "github.com/parinithshekar/gitsink/common/config"
	logger "github.com/parinithshekar/gitsink/wrap/logrus/v1"
)

// Public struct defines fields in github-public object
type Public struct {
	accountID   string
	accessToken string
	kind        string
	api         *github.Client
	ctx         context.Context
}

// setAPIClient adds a usable API client to the initiated struct
func (public *Public) setAPIClient() {
	ctx := context.Background()

	accessToken := os.Getenv(public.accessToken)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	public.api = github.NewClient(tc)
	public.ctx = ctx
}

// Credentials fetches amd returns the accountID and accessToken from environment variables
func (public Public) Credentials() (string, string) {
	accountID := os.Getenv(public.accountID)
	accessToken := os.Getenv(public.accessToken)
	return accountID, accessToken
}

// New returns a new github-public object
func New(target config.Target) *Public {
	var public *Public = new(Public)

	// Check if env variables mentioned in config file exist
	// check account ID env variable
	_, exists := os.LookupEnv(target.AccountID)
	if !exists {
		logger.New().Errorf("Account ID not found")
		os.Exit(1)
	} else {
		public.accountID = target.AccountID
	}
	// check access token env variable
	_, exists = os.LookupEnv(target.AccessToken)
	if !exists {
		logger.New().Errorf("Access Token not found")
		os.Exit(1)
	} else {
		public.accessToken = target.AccessToken
	}

	public.kind = target.Kind

	public.setAPIClient()
	return public
}

// Authenticate checks the account ID and access tokens' validity for the kind defined
func (public Public) Authenticate() (bool, error) {
	kindSplit := strings.SplitN(public.kind, "/", 2)
	kindType := kindSplit[0]
	kindKey := kindSplit[1]

	accountID := os.Getenv(public.accountID)

	switch kindType {
	case "org":
		// returns true only if the user is a member of the org, can create repositories
		// returns Membership, Response, error
		_, _, err := public.api.Organizations.GetOrgMembership(public.ctx, accountID, kindKey)
		if err != nil {
			return false, err
		}
		return true, nil

	case "user":
		// return true if the authenticated user from env variables is the same user mentioned in config
		user, _, err := public.api.Users.Get(context.Background(), "")
		if err != nil {
			fmt.Println("USER AUTH: ", err.Error())
			return false, err
		} else if kindKey != *user.Login {
			return false, fmt.Errorf("Kind username does not match account ID. Unable to push to this target")
		}
		return true, nil

	default:
		// Mentioned kind is unsupported
		return false, fmt.Errorf("Unsupported kind: %v", kindType)
	}
}

func (public Public) makeNewRepo(repo common.Repository) (string, error) {
	kindSplit := strings.SplitN(public.kind, "/", 2)
	kindType := kindSplit[0]
	kindKey := kindSplit[1]

	newRepository := github.Repository{
		Name:        &repo.Slug,
		Description: &repo.Description,
	}

	// Make new repo
	var err error
	var newRepo *github.Repository
	switch kindType {
	case "org":
		newRepo, _, err = public.api.Repositories.Create(context.Background(), kindKey, &newRepository)

	case "user":
		newRepo, _, err = public.api.Repositories.Create(context.Background(), "", &newRepository)
	}
	if err != nil {
		logger.New().Errorf("Failed to make new repository: %v", repo.Slug)
		return "", fmt.Errorf("Failed to make new repository: %v", repo.Slug)
	}

	newRepoBytes, _ := json.MarshalIndent(newRepo, "", "  ")
	newRepoJSON := string(newRepoBytes)
	targetURL := gjson.Get(newRepoJSON, `clone_url`).String()
	return targetURL, nil
}

// SyncCheck checks whether the repository is already present at the target
// If it is, then only a sync is done, else a new repository is created at the target
func (public Public) SyncCheck(repos []common.Repository) []common.Repository {
	kindSplit := strings.SplitN(public.kind, "/", 2)
	// kindType := kindSplit[0]
	kindKey := kindSplit[1]

	var processedRepos []common.Repository

	for _, repo := range repos {

		targetRepo, _, err := public.api.Repositories.Get(context.Background(), kindKey, repo.Slug)

		if err != nil {
			logger.New().Infof("%v - Repository not found", repo.Slug)
			targetURL, err := public.makeNewRepo(repo)
			if err != nil {
				logger.New().Warningf("Skipping repository: %v", repo.Slug)
				continue
			} else {
				repo.Target = targetURL
			}
		} else {
			targetRepoBytes, _ := json.MarshalIndent(targetRepo, "", "  ")
			targetRepoJSON := string(targetRepoBytes)
			repo.Target = gjson.Get(targetRepoJSON, `clone_url`).String()
		}

		processedRepos = append(processedRepos, repo)

	}
	return processedRepos
}
