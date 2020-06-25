package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	gjson "github.com/tidwall/gjson"
	// bitbucketv1 "github.com/gfleury/go-bitbucket-v1"

	common "github.com/parinithshekar/gitsink/common"
	config "github.com/parinithshekar/gitsink/common/config"
	logger "github.com/parinithshekar/gitsink/wrap/logrus/v1"
)

// Server struct defines the data fields in bitbucket-server object
type Server struct {
	apiBaseURL string

	accountID   string
	accessToken string
	kind        string
	api         *http.Client
}

// Credentials fetches amd returns the accountID and accessToken from environment variables
func (server *Server) Credentials() (string, string) {
	accountID := os.Getenv(server.accountID)
	accessToken := os.Getenv(server.accessToken)

	return accountID, accessToken
}

// setAPIClient builds and returns an object to facilitate calls to the API
func (server *Server) setAPIClient(baseURL string) {

	client := &http.Client{
		Timeout: time.Second * 15,
	}

	server.apiBaseURL = baseURL + "/bitbucket/rest/api/1.0"
	server.api = client
}

// New returns a new bitbucket-server object with metadata
func New(source config.Source) *Server {
	var server *Server = new(Server)

	_, exists := os.LookupEnv(source.AccountID)
	if !exists {
		logger.New().Errorf("Account ID not found")
		os.Exit(1)
	} else {
		server.accountID = source.AccountID
	}

	_, exists = os.LookupEnv(source.AccessToken)
	if !exists {
		logger.New().Errorf("Access Token not found")
		os.Exit(1)
	} else {
		server.accessToken = source.AccessToken
	}

	server.kind = source.Kind

	server.setAPIClient(source.BaseURL)

	return server
}

// Authenticate checks the account ID and access tokens' validity for the kind defined
func (server *Server) Authenticate() (bool, error) {
	kindSplit := strings.Split(server.kind, "/")
	kindType := kindSplit[0]
	kindKey := kindSplit[1]

	accountID, accessToken := server.Credentials()

	switch kindType {
	case "project":
		request, err := http.NewRequest("GET", server.apiBaseURL+"/projects/"+kindKey+"/repos", nil)
		request.SetBasicAuth(accountID, accessToken)

		// Check if user can access repos of the project mentioned (kindKey) in config
		_, err = server.api.Do(request)
		if err != nil {
			fmt.Printf("Authenticate Project: %v\n", err.Error())
			return false, err
		}
		return true, nil

	case "user":
		request, err := http.NewRequest("GET", server.apiBaseURL+"/users/"+kindKey+"/repos", nil)
		request.SetBasicAuth(accountID, accessToken)

		// Check if user can access repos of the user mentioned (kindKey) in config
		_, err = server.api.Do(request)
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

// allRepositories abstracts over paginated results and gives a list of all the repos
func (server *Server) allRepositories(URL, accountID, accessToken string) ([]common.Repository, error) {
	isLastPage := false
	var start int64 = 0

	repositories := []common.Repository{}

	for !isLastPage {
		pagedURL := fmt.Sprintf("%v?start=%v", URL, start)
		request, err := http.NewRequest("GET", pagedURL, nil)
		request.SetBasicAuth(accountID, accessToken)

		response, err := server.api.Do(request)
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
	// The bitbucketv1 API client does not abstract over pagination
	// Might as well write this from scratch using HTTP calls
	kindSplit := strings.Split(server.kind, "/")
	kindType := kindSplit[0]
	kindKey := kindSplit[1]

	fmt.Println("INSIDE REPOSITORIES")
	accountID, accessToken := server.Credentials()

	switch kindType {
	case "project":
		reposURL := fmt.Sprintf("%v/projects/%v/repos", server.apiBaseURL, kindKey)
		// abstract over pagination
		repositories, err := server.allRepositories(reposURL, accountID, accessToken)
		if err != nil {
			fmt.Printf("Get Repositories: %v\n", err.Error())
			return nil, err
		}
		return repositories, nil

	case "user":
		reposURL := fmt.Sprintf("%v/users/%v/repos", server.apiBaseURL, kindKey)
		// abstract over pagination
		repositories, err := server.allRepositories(reposURL, accountID, accessToken)
		if err != nil {
			fmt.Printf("Get Repositories: %v\n", err.Error())
			return nil, err
		}
		return repositories, nil

	default:
		return nil, fmt.Errorf("Unsupported kind: %v", kindType)
	}
}
