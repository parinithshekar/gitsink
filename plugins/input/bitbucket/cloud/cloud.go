package cloud

import (
	"os"
	"fmt"
	"strings"
	"encoding/json"

	gjson "github.com/tidwall/gjson"
	bitbucket "github.com/ktrysmt/go-bitbucket"

	common "github.com/parinithshekar/gitsink/common"
	config "github.com/parinithshekar/gitsink/common/config"
	logger "github.com/parinithshekar/gitsink/wrap/logrus/v1"
)

// Cloud struct defines data fields in bitbucket-cloud object
type Cloud struct {
	accountID   string
	accessToken string
	kind        string
	api         *bitbucket.Client
}

// setAPIClient builds and returns an object to facilitate calls to the API
func (cloud *Cloud) setAPIClient() {
	accountID := os.Getenv(cloud.accountID)
	accessToken := os.Getenv(cloud.accessToken)

	client := bitbucket.NewBasicAuth(accountID, accessToken)
	cloud.api = client
}

// Credentials fetches amd returns the accountID and accessToken from environment variables
func (cloud Cloud) Credentials() (string, string) {
	accountID := os.Getenv(cloud.accountID)
	accessToken := os.Getenv(cloud.accessToken)
	return accountID, accessToken
}

// New returns a new bitbucket-cloud object
func New(source config.Source) *Cloud {
	var cloud *Cloud = new(Cloud)

	_, exists := os.LookupEnv(source.AccountID)
	if !exists {
		logger.New().Errorf("Account ID not found")
		os.Exit(1)
	} else {
		cloud.accountID = source.AccountID
	}

	_, exists = os.LookupEnv(source.AccessToken)
	if !exists {
		logger.New().Errorf("Access Token not found")
		os.Exit(1)
	} else {
		cloud.accessToken = source.AccessToken
	}

	cloud.kind = source.Kind

	cloud.setAPIClient()
	return cloud
}

// Authenticate checks the account ID and access tokens' validity for the kind defined
func (cloud Cloud) Authenticate() (bool, error) {
	kindSplit := strings.SplitN(cloud.kind, "/", 2)
	kindType := kindSplit[0]
	kindKey := kindSplit[1]

	accountID := os.Getenv(cloud.accountID)

	switch kindType {
	case "project":
		/*
			The go-bitbucket SDK does not cover all API endpoints for some reason
			Will have to get list of projects and search for the project key
		*/

		// Get list of projects user has access to
		result, err := cloud.api.Teams.Projects(accountID)
		if err != nil {
			fmt.Printf("Authenticate Project: %v\n", err.Error())
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
		return false, fmt.Errorf("Project key not found: %s", kindKey)

	case "user":
		// Check if current user can access the repos of the user mentioned (kindKey) in config
		_, err := cloud.api.Teams.Repositories(kindKey)
		if err != nil {
			fmt.Printf("Authenticate User: %v\n", err.Error())
			return false, err
		}
		return true, nil

	default:
		// Mentioned kind is unsupported
		return false, fmt.Errorf("Unsupported kind: %v", kindType)
	}
}

// Repositories queries the API and returns a list of repositories mentioned by the kind
func (cloud Cloud) Repositories(metadata bool) ([]common.Repository, error) {
	kindSplit := strings.SplitN(cloud.kind, "/", 2)
	kindType := kindSplit[0]
	kindKey := kindSplit[1]

	accountID := os.Getenv(cloud.accountID)

	switch kindType {
	case "project":
		ro := bitbucket.RepositoriesOptions{Owner: accountID, Role: "member"}
		result, err := cloud.api.Repositories.ListForAccount(&ro)
		if err != nil {
			fmt.Printf("Get project repositories: %v", err.Error())
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
		result, err := cloud.api.Repositories.ListForAccount(&ro)
		if err != nil {
			fmt.Printf("Get repositories: %v", err.Error())
			return nil, err
		}

		repositories := []common.Repository{}
		for _, repo := range result.Items {
			newRepo := common.Repository{
				Slug: repo.Slug,
			}
			repositories = append(repositories, newRepo)
		}
		return repositories, nil

	default:
		return nil, fmt.Errorf("Unsupported kind: %v", kindType)
	}
}
