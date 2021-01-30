package bbcloud

import (
	"errors"
	bitbucket "github.com/ktrysmt/go-bitbucket"
)

// Teams ...
type Teams struct {
	AccountID, AccessToken string
}

// Repositories ...
type Repositories struct {
	AccountID, AccessToken string
}

// MockAPI ...
type MockAPI struct {
	Teams        Teams
	Repositories Repositories
}

// Projects ...
func (teams *Teams) Projects(kindKey string) (interface{}, error) {

	// Check access token
	if teams.AccessToken != "token" {
		return nil, errors.New("Wrong access token")
	}

	// Check account ID
	if teams.AccountID != "username" {
		return nil, errors.New("Access denied")
	}

	// Check kind key
	if kindKey != "username" {
		// No "TEST" project key
		var wrongRes interface{} = map[string]interface{}{
			"values": []interface{}{
				map[string]interface{}{
					"key": "JETT",
				},
				map[string]interface{}{
					"key": "RAZE",
				},
			},
		}
		return wrongRes, nil
	}

	// Has "TEST" project key
	var correctRes interface{} = map[string]interface{}{
		"values": []interface{}{
			map[string]interface{}{
				"key": "ABC",
			},
			map[string]interface{}{
				"key": "POKE",
			},
			map[string]interface{}{
				"key": "TEST",
			},
		},
	}
	return correctRes, nil
}

// Repositories ...
func (teams *Teams) Repositories(kindKey string) (interface{}, error) {

	/* This function is only used for the authentication check.
	Therefore, returning a dummy interface */

	// Check access token
	if teams.AccessToken != "token" {
		return nil, errors.New("Bad credentials")
	}

	// Check account ID
	if teams.AccountID != "username" {
		return nil, errors.New("Access denied")
	}

	// Check kind key
	if kindKey != "username" {
		return nil, errors.New("Wrong kind key")
	}

	var res interface{} = "Team.Repositories"
	return res, nil
}

// ListForAccount ...
func (repositories *Repositories) ListForAccount(ro *bitbucket.RepositoriesOptions) (*bitbucket.RepositoriesRes, error) {

	// Check access token
	if repositories.AccessToken != "token" {
		return nil, errors.New("Bad credentials")
	}

	// Check account ID
	if repositories.AccountID != "username" {
		return nil, errors.New("Access denied")
	}

	// Check kind key
	if ro.Owner != "username" {
		return nil, errors.New("Access Denied")
	}

	var res bitbucket.RepositoriesRes = bitbucket.RepositoriesRes{
		Page:     1,
		Pagelen:  3,
		MaxDepth: 1,
		Size:     3,
		Items: []bitbucket.Repository{
			{
				Project: bitbucket.Project{
					Key:  "WRONG",
					Name: "Rolling Thunder",
				},
				Slug:        "repo-1",
				Full_name:   "repo-1",
				Description: "repository for testing",
				ForkPolicy:  "fork_policy",
				Type:        "repository",
				Owner: map[string]interface{}{
					"one": 1,
				},
				Links: map[string]interface{}{
					"clone": []interface{}{
						map[string]string{
							"name": "ssh",
							"href": "@ssh+repo1.link",
						},
						map[string]string{
							"name": "https",
							"href": "https://gitclub.com/username/repo-1.git",
						},
					},
				},
			},
			{
				Project: bitbucket.Project{
					Key:  "TEST",
					Name: "testing BB Cloud",
				},
				Slug:        "repo-2",
				Full_name:   "repo-2",
				Description: "repository for testing",
				ForkPolicy:  "fork_policy",
				Type:        "repository",
				Owner: map[string]interface{}{
					"one": 1,
				},
				Links: map[string]interface{}{
					"clone": []interface{}{
						map[string]string{
							"name": "ssh",
							"href": "@ssh+repo2.link",
						},
						map[string]string{
							"name": "https",
							"href": "https://gitclub.com/username/repo-2.git",
						},
					},
				},
			},
			{
				Project: bitbucket.Project{
					Key:  "TEST",
					Name: "testing BB Cloud",
				},
				Slug:        "repo-3",
				Full_name:   "repo-3",
				Description: "repository for testing",
				ForkPolicy:  "fork_policy",
				Type:        "repository",
				Owner: map[string]interface{}{
					"one": 1,
				},
				Links: map[string]interface{}{
					"clone": []interface{}{
						map[string]string{
							"name": "ssh",
							"href": "@ssh+repo3.link",
						},
						map[string]string{
							"name": "https",
							"href": "https://gitclub.com/username/repo-3.git",
						},
					},
				},
			},
		},
	}
	return &res, nil
}
