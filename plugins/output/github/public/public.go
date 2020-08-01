package public

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	github "github.com/google/go-github/v31/github"
	logrus "github.com/sirupsen/logrus"
	gjson "github.com/tidwall/gjson"
	oauth2 "golang.org/x/oauth2"

	common "github.com/parinithshekar/gitsink/common"
	config "github.com/parinithshekar/gitsink/common/config"
	logger "github.com/parinithshekar/gitsink/wrap/logrus/v1"
)

var (
	log = logger.New()
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
func (public Public) Credentials() (string, string, error) {
	accountID := os.Getenv(public.accountID)
	accessToken := os.Getenv(public.accessToken)

	accountID, exists := os.LookupEnv(public.accountID)
	if !exists {
		log.WithFields(logrus.Fields{
			"accountID": public.accountID,
		}).Errorf("Account ID not found")
		return "", "", fmt.Errorf("Account ID not found")
	}

	accessToken, exists = os.LookupEnv(public.accessToken)
	if !exists {
		log.WithFields(logrus.Fields{
			"accessToken": public.accessToken,
		}).Errorf("Access Token not found")
		return "", "", fmt.Errorf("Access Token not found")
	}

	return accountID, accessToken, nil
}

// New returns a new github-public object
func New(target config.Target) (*Public, error) {
	var public *Public = new(Public)

	// Check if env variables mentioned in config file exist
	// check account ID env variable
	_, exists := os.LookupEnv(target.AccountID)
	if !exists {
		log.WithFields(logrus.Fields{
			"accountID": target.AccountID,
		}).Errorf("Account ID not found")
		return nil, fmt.Errorf("Account ID not found")
	}
	public.accountID = target.AccountID
	// check access token env variable
	_, exists = os.LookupEnv(target.AccessToken)
	if !exists {
		log.WithFields(logrus.Fields{
			"accessToken": target.AccessToken,
		}).Errorf("Access Token not found")
		return nil, fmt.Errorf("Access Token not found")
	}
	public.accessToken = target.AccessToken

	public.kind = target.Kind

	public.setAPIClient()
	return public, nil
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
			log.WithFields(logrus.Fields{
				"organization": kindKey,
			}).Errorf("Organization membership check failed")
			return false, err
		}
		return true, nil

	case "user":
		// return true if the authenticated user from env variables is the same user mentioned in config
		user, _, err := public.api.Users.Get(context.Background(), "")
		if err != nil {
			return false, err
		} else if kindKey != *user.Login {
			log.WithFields(logrus.Fields{
				"user": kindKey,
			}).Errorf("Kind username does not match account ID")
			return false, fmt.Errorf("Unable to push to target user")
		}
		return true, nil

	default:
		// Mentioned kind is unsupported
		log.WithFields(logrus.Fields{
			"kind": kindType,
		}).Errorf("Unsupported kind")
		return false, fmt.Errorf("Unsupported kind")
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
		log.WithFields(logrus.Fields{
			"repository": repo.Slug,
		}).Errorf("Repository creation failed")
		return "", fmt.Errorf("Repository creation failed")
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
			log.WithFields(logrus.Fields{
				"repository": repo.Slug,
			}).Infof("Repository not found")
			targetURL, err := public.makeNewRepo(repo)
			if err != nil {
				log.WithFields(logrus.Fields{
					"repository": repo.Slug,
				}).Warningf("Skipping repository")
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
