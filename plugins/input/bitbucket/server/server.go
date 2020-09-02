package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	logrus "github.com/sirupsen/logrus"
	gjson "github.com/tidwall/gjson"

	common "github.com/parinithshekar/gitsink/common"
	config "github.com/parinithshekar/gitsink/common/config"
	utils "github.com/parinithshekar/gitsink/common/utils"
	logger "github.com/parinithshekar/gitsink/wrap/logrus/v1"
)

var (
	log = logger.New()
)

// APIClient defines the methods for the API in Server object
type APIClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Server struct defines the data fields in bitbucket-server object
type Server struct {
	apiBaseURL  string
	accountID   string
	accessToken string
	kind        string
	filters     struct {
		include []string
		exclude []string
	}
	API APIClient
}

// Credentials fetches amd returns the accountID and accessToken from environment variables
func (server *Server) Credentials() (string, string, error) {
	accountID := os.Getenv(server.accountID)
	accessToken := os.Getenv(server.accessToken)

	accountID, exists := os.LookupEnv(server.accountID)
	if !exists {
		log.WithFields(logrus.Fields{
			"accountID": server.accountID,
		}).Errorf("Account ID not found")
		return "", "", fmt.Errorf("Account ID not found")
	}

	accessToken, exists = os.LookupEnv(server.accessToken)
	if !exists {
		log.WithFields(logrus.Fields{
			"accessToken": server.accessToken,
		}).Errorf("Access Token not found")
		return "", "", fmt.Errorf("Access Token not found")
	}

	return accountID, accessToken, nil
}

// setAPIClient builds and returns an object to facilitate calls to the API
func (server *Server) setAPIClient(baseURL string) {

	client := &http.Client{
		Timeout: time.Second * 15,
	}

	server.apiBaseURL = baseURL + "/bitbucket/rest/api/1.0"
	server.API = client
}

// New returns a new bitbucket-server object with metadata
func New(source config.Source) (*Server, error) {
	var server *Server = new(Server)

	_, exists := os.LookupEnv(source.AccountID)
	if !exists {
		log.WithFields(logrus.Fields{
			"accountID": source.AccountID,
		}).Errorf("Account ID not found")
		return nil, fmt.Errorf("Account ID not found")
	}
	server.accountID = source.AccountID

	_, exists = os.LookupEnv(source.AccessToken)
	if !exists {
		log.WithFields(logrus.Fields{
			"accessToken": source.AccessToken,
		}).Errorf("Access Token not found")
		return nil, fmt.Errorf("Access Token not found")
	}
	server.accessToken = source.AccessToken

	server.kind = source.Kind

	server.filters.include = source.Repositories.Include
	server.filters.exclude = source.Repositories.Exclude

	server.setAPIClient(source.BaseURL)

	return server, nil
}

// Authenticate checks the account ID and access tokens' validity for the kind defined
func (server *Server) Authenticate() (bool, error) {

	kindSplit := strings.Split(server.kind, "/")
	kindType := kindSplit[0]
	kindKey := kindSplit[1]

	accountID, accessToken, err := server.Credentials()
	if err != nil {
		log.WithFields(logrus.Fields{
			"accountID":   server.accountID,
			"accessToken": server.accessToken,
		}).Errorf("Failed to fetch credentials")
		return false, err
	}

	switch kindType {
	case "project":
		request, err := http.NewRequest("GET", server.apiBaseURL+"/projects/"+kindKey+"/repos", nil)
		request.SetBasicAuth(accountID, accessToken)

		// Check if user can access repos of the project mentioned (kindKey) in config
		_, err = server.API.Do(request)
		if err != nil {
			log.WithFields(logrus.Fields{
				"project": kindKey,
			}).Errorf("Project not found. Check user access")
			return false, err
		}
		return true, nil

	case "user":
		request, err := http.NewRequest("GET", server.apiBaseURL+"/users/"+kindKey+"/repos", nil)
		request.SetBasicAuth(accountID, accessToken)

		// Check if user can access repos of the user mentioned (kindKey) in config
		_, err = server.API.Do(request)
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

// allRepositories abstracts over paginated results and gives a list of all the repos
func (server *Server) allRepositories(URL, accountID, accessToken string) ([]common.Repository, error) {
	isLastPage := false
	var start int64 = 0

	repositories := []common.Repository{}

	for !isLastPage {
		pagedURL := fmt.Sprintf("%v?start=%v", URL, start)
		request, err := http.NewRequest("GET", pagedURL, nil)
		request.SetBasicAuth(accountID, accessToken)

		response, err := server.API.Do(request)
		if err != nil {
			return nil, err
		}
		bodyBytes, err := ioutil.ReadAll(response.Body)
		bodyJSON := string(bodyBytes)

		// Continue fetching pages until last page
		isLastPage = gjson.Get(bodyJSON, "isLastPage").Bool()
		if !isLastPage {
			start = gjson.Get(bodyJSON, "nextPageStart").Int()
		}

		repos := gjson.Get(bodyJSON, "values").Array()

		for _, repoJSON := range repos {

			// Get repo metadata
			httpCloneLink := gjson.Get(repoJSON.String(), `links.clone.#(name%"http*").href`).String()
			slug := gjson.Get(repoJSON.String(), `slug`).String()
			description := gjson.Get(repoJSON.String(), `description`).String()

			newRepo := common.Repository{
				Slug:        slug,
				Source:      httpCloneLink,
				Description: description,
			}
			repositories = append(repositories, newRepo)
		}
	}
	return repositories, nil
}

// Repositories queries the API and returns a list of repositories mentioned by the kind
func (server *Server) Repositories(metadata bool) ([]common.Repository, error) {

	kindSplit := strings.Split(server.kind, "/")
	kindType := kindSplit[0]
	kindKey := kindSplit[1]

	accountID, accessToken, err := server.Credentials()
	if err != nil {
		log.WithFields(logrus.Fields{
			"accountID":   server.accountID,
			"accessToken": server.accessToken,
		}).Errorf("Failed to fetch credentials")
		return nil, err
	}

	switch kindType {
	case "project":
		reposURL := fmt.Sprintf("%v/projects/%v/repos", server.apiBaseURL, kindKey)
		// abstract over pagination
		repositories, err := server.allRepositories(reposURL, accountID, accessToken)
		if err != nil {
			log.WithFields(logrus.Fields{
				"project": kindKey,
			}).Errorf("Failed to get project repositories")
			return nil, err
		}
		repositories = utils.FilterRepos(repositories, server.filters.include, server.filters.exclude)
		return repositories, nil

	case "user":
		reposURL := fmt.Sprintf("%v/users/%v/repos", server.apiBaseURL, kindKey)
		// abstract over pagination
		repositories, err := server.allRepositories(reposURL, accountID, accessToken)
		if err != nil {
			log.WithFields(logrus.Fields{
				"user": kindKey,
			}).Errorf("Failed to get user repositories")
			return nil, err
		}
		repositories = utils.FilterRepos(repositories, server.filters.include, server.filters.exclude)
		return repositories, nil

	default:
		log.WithFields(logrus.Fields{
			"kind": kindType,
		}).Errorf("Unsupported kind")
		return nil, fmt.Errorf("Unsupported kind")
	}
}
