package bbserver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// MockAPI is a mock API client to help with testing
type MockAPI struct {
	BaseURL string
}

// Link ...
type Link struct {
	Name string `json:"name"`
	Href string `json:"href"`
}

// Value ...
type Value struct {
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	Links       map[string][]Link `json:"links"`
}

// Repos defines the response format for the API
type Repos struct {
	IsLastPage    bool    `json:"isLastPage"`
	NextPageStart int     `json:"nextPageStart"`
	Values        []Value `json:"values"`
}

var (
	projectSuffix = "/projects/TEST/repos"
	userSuffix    = "/users/username/repos"

	// DoFunc is called by the mock's Do function
	DoFunc = func(req *http.Request) (*http.Response, error) {

		// The request can have extra parameters as suffix ?start=25
		URL := req.URL.EscapedPath()
		projectRequest := strings.Contains(URL, projectSuffix)
		userRequest := strings.Contains(URL, userSuffix)

		// Get username and access token
		auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			return nil, errors.New("Authorization Failed: Bad Header")
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			return nil, errors.New("Authorization Failed: Missing username or token")
		}
		username := pair[0]
		token := pair[1]
		ok := (username == "username") && (token == "token")
		if !ok {
			return nil, errors.New("Bad credentials")
		}

		switch {
		case projectRequest:
			repos := Repos{
				IsLastPage: true,
				Values: []Value{
					{
						Slug:        "project-repo-1",
						Description: "describe project-repo-1",
						Links: map[string][]Link{
							"clone": {
								{
									Name: "https",
									Href: "https://www.bitbucket-abc/TEST/project-repo-1",
								},
								{
									Name: "ssh",
									Href: "git@bitbucket-abc/company.com:TEST/project-repo-1.git",
								},
							},
						},
					},
				},
			}
			reposBytes, _ := json.Marshal(repos)
			response := http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(string(reposBytes)))),
			}
			return &response, nil
		case userRequest:
			repos := Repos{
				IsLastPage: true,
				Values: []Value{
					{
						Slug:        "user-repo-1",
						Description: "describe user-repo-1",
						Links: map[string][]Link{
							"clone": {
								{
									Name: "https",
									Href: "https://www.bitbucket-abc/username/user-repo-1",
								},
								{
									Name: "ssh",
									Href: "git@bitbucket-abc/company.com:username/user-repo-1.git",
								},
							},
						},
					},
				},
			}
			reposBytes, _ := json.Marshal(repos)
			response := http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(string(reposBytes)))),
			}
			return &response, nil
		default:
			return nil, errors.New("Unsupported kind")
		}
	}
)

// Do is the mock Do method that mimics http package functionality
func (m *MockAPI) Do(req *http.Request) (*http.Response, error) {
	return DoFunc(req)
}
