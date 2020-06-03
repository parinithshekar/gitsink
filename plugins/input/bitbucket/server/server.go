package server

import (
	"os"
	"strings"
	"context"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
)

// Server struct defines the data fields in bitbucket-server object
type Server struct {
	BaseURL    string
	APIBaseURL string
	accountID  string
	api        *bitbucketv1.APIClient
}

// setAPIClient builds and returns an object to facilitate calls to the API
func (server *Server) setAPIClient(envAccountID, envAccessToken string) {
	accountID := os.Getenv(envAccountID)
	accessToken := os.Getenv(envAccessToken)
	basicAuth := bitbucketv1.BasicAuth{UserName: accountID, Password: accessToken}

	ctx := context.WithValue(context.Background(), bitbucketv1.ContextBasicAuth, basicAuth)
	server.api = bitbucketv1.NewAPIClient(ctx, bitbucketv1.NewConfiguration(server.APIBaseURL))
}

// New returns a new bitbucket-server object with metadata
func New(baseURL, envAccountID, envAccessToken string) *Server {
	var server *Server = new(Server)
	server.BaseURL = baseURL
	server.APIBaseURL = baseURL + "/bitbucket/rest"
	server.accountID = os.Getenv(envAccountID)
	server.setAPIClient(envAccountID, envAccessToken)

	return server
}

// Authenticate checks the account ID and access tokens' validity for the kind defined
func (server *Server) Authenticate(kind string) (bool, error) {
	kindSplit := strings.Split(kind, "/")
	projectKey := kindSplit[1]
	result, err := server.api.DefaultApi.GetProject(projectKey)
	if result.StatusCode != 200 {
		return false, err
	}
	return true, nil
}

// Repositories queries the API and returns a list of repositories mentioned by the kind
func (server *Server) Repositories(kind string, metadata bool) []string {
	// The bitbucketv1 API client does not abstract over pagination
	// Might as well write this from scratch using HTTP calls
	dummy := []string{"repo1", "repo2"}
	return dummy
}
