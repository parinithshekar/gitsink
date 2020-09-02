package cloud

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	bitbucket "github.com/ktrysmt/go-bitbucket"
	logrus "github.com/sirupsen/logrus"
	gjson "github.com/tidwall/gjson"

	common "github.com/parinithshekar/gitsink/common"
	config "github.com/parinithshekar/gitsink/common/config"
	logger "github.com/parinithshekar/gitsink/wrap/logrus/v1"
)

var (
	log = logger.New()
)

// Teams interface declares methods to implement in the API object
type Teams interface {
	Projects(string) (interface{}, error)
	Repositories(string) (interface{}, error)
}

// Repositories interface declares methods to implement in the API object
type Repositories interface {
	ListForAccount(*bitbucket.RepositoriesOptions) (*bitbucket.RepositoriesRes, error)
}

// Cloud struct defines data fields in bitbucket-cloud object
type Cloud struct {
	accountID   string
	accessToken string
	kind        string
	API         struct {
		Teams        Teams
		Repositories Repositories
	}
}

// setAPIClient builds and returns an object to facilitate calls to the API
func (cloud *Cloud) setAPIClient() {

	accountID, accessToken, _ := cloud.Credentials()

	client := bitbucket.NewBasicAuth(accountID, accessToken)
	cloud.API.Teams = client.Teams
	cloud.API.Repositories = client.Repositories
}

// Credentials fetches amd returns the accountID and accessToken from environment variables
func (cloud Cloud) Credentials() (string, string, error) {
	accountID := os.Getenv(cloud.accountID)
	accessToken := os.Getenv(cloud.accessToken)

	accountID, exists := os.LookupEnv(cloud.accountID)
	if !exists {
		log.WithFields(logrus.Fields{
			"accountID": cloud.accountID,
		}).Errorf("Account ID not found")
		return "", "", fmt.Errorf("Account ID not found")
	}

	accessToken, exists = os.LookupEnv(cloud.accessToken)
	if !exists {
		log.WithFields(logrus.Fields{
			"accessToken": cloud.accessToken,
		}).Errorf("Access Token not found")
		return "", "", fmt.Errorf("Access Token not found")
	}

	return accountID, accessToken, nil
}

// New returns a new bitbucket-cloud object
func New(source config.Source) (*Cloud, error) {
	var cloud *Cloud = new(Cloud)

	_, exists := os.LookupEnv(source.AccountID)
	if !exists {
		log.WithFields(logrus.Fields{
			"accountID": source.AccountID,
		}).Errorf("Account ID not found")
		return nil, fmt.Errorf("Account ID not found")
	}
	cloud.accountID = source.AccountID

	_, exists = os.LookupEnv(source.AccessToken)
	if !exists {
		log.WithFields(logrus.Fields{
			"accessToken": source.AccessToken,
		}).Errorf("Access Token not found")
		return nil, fmt.Errorf("Access Token not found")
	}
	cloud.accessToken = source.AccessToken

	cloud.kind = source.Kind

	cloud.setAPIClient()
	return cloud, nil
}

// Authenticate checks the account ID and access tokens' validity for the kind defined
func (cloud Cloud) Authenticate() (bool, error) {

	kindSplit := strings.SplitN(cloud.kind, "/", 2)
	kindType := kindSplit[0]
	kindKey := kindSplit[1]

	accountID, _, err := cloud.Credentials()
	if err != nil {
		log.Errorf("Failed to authenticate")
		return false, err
	}

	switch kindType {
	case "project":
		/*
			The go-bitbucket SDK does not cover all API endpoints for some reason
			Will have to get list of projects and search for the project key
		*/

		// Get list of projects user has access to
		result, err := cloud.API.Teams.Projects(accountID)
		if err != nil {
			log.Errorf("Failed to get projects")
			return false, err
		}

		// Check if mentioned project is in the list
		passed := false
		values := result.(map[string]interface{})["values"].([]interface{})
		for _, v := range values {
			passed = kindKey == v.(map[string]interface{})["key"]
			if passed {
				return passed, nil
			}
		}
		log.WithFields(logrus.Fields{
			"project": kindKey,
		}).Errorf("Project not found. Check user access")
		return false, fmt.Errorf("Project not found. Check user access")

	case "user":
		// Check if current user can access the repos of the user mentioned (kindKey) in config
		_, err := cloud.API.Teams.Repositories(kindKey)
		if err != nil {
			log.WithFields(logrus.Fields{
				"user": kindKey,
			}).Errorf("User authentication failed")
			return false, err
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

// Repositories queries the API and returns a list of repositories mentioned by the kind
func (cloud Cloud) Repositories(metadata bool) ([]common.Repository, error) {

	kindSplit := strings.SplitN(cloud.kind, "/", 2)
	kindType := kindSplit[0]
	kindKey := kindSplit[1]

	accountID, _, err := cloud.Credentials()
	if err != nil {
		log.Errorf("Failed to get repositories")
		return nil, err
	}

	switch kindType {
	case "project":
		ro := bitbucket.RepositoriesOptions{Owner: accountID, Role: "member"}
		result, err := cloud.API.Repositories.ListForAccount(&ro)
		if err != nil {
			log.WithFields(logrus.Fields{
				"project": kindKey,
			}).Errorf("Failed to get project repositories")
			return nil, err
		}

		repositories := []common.Repository{}
		for _, repo := range result.Items {
			if repo.Project.Key == kindKey {

				repoBytes, _ := json.MarshalIndent(repo, "", "  ")
				repoJSON := string(repoBytes)
				// fmt.Println(repoJSON)
				// Get all metadata for the repository
				httpCloneLink := gjson.Get(repoJSON, `Links.clone.#(name%"http*").href`).String()
				slug := gjson.Get(repoJSON, `Slug`).String()
				description := gjson.Get(repoJSON, `Description`).String()

				newRepo := common.Repository{
					Slug:        slug,
					Source:      httpCloneLink,
					Description: description,
				}
				repositories = append(repositories, newRepo)
			}
		}
		return repositories, nil

	case "user":
		ro := bitbucket.RepositoriesOptions{Owner: kindKey, Role: "member"}
		result, err := cloud.API.Repositories.ListForAccount(&ro)
		if err != nil {
			log.WithFields(logrus.Fields{
				"user": kindKey,
			}).Errorf("Failed to get user repositories")
			return nil, err
		}

		repositories := []common.Repository{}
		for _, repo := range result.Items {
			repoBytes, _ := json.MarshalIndent(repo, "", "  ")
			repoJSON := string(repoBytes)

			httpCloneLink := gjson.Get(repoJSON, `Links.clone.#(name%"http*").href`).String()
			slug := gjson.Get(repoJSON, `Slug`).String()
			description := gjson.Get(repoJSON, `Description`).String()

			newRepo := common.Repository{
				Slug:        slug,
				Source:      httpCloneLink,
				Description: description,
			}
			repositories = append(repositories, newRepo)
		}
		return repositories, nil

	default:
		log.WithFields(logrus.Fields{
			"kind": kindType,
		}).Errorf("Unsupported kind")
		return nil, fmt.Errorf("Unsupported kind")
	}
}
