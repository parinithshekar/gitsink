package cloud_test

import (
	"os"
	"testing"

	config "github.com/parinithshekar/gitsink/common/config"
	mock "github.com/parinithshekar/gitsink/mocks/bbcloud"
	bbcloud "github.com/parinithshekar/gitsink/plugins/input/bitbucket/cloud"
	plugins "github.com/parinithshekar/gitsink/plugins/interfaces"
)

var (
	envAccountID   = "TEST_BBCLOUD_ACCOUNT_ID"
	envAccessToken = "TEST_BBCLOUD_ACCESS_TOKEN"

	source config.Source = config.Source{
		Type:        "bitbucket-cloud",
		AccountID:   envAccountID,
		AccessToken: envAccessToken,
		Kind:        "user/username",
		Repositories: config.Repositories{
			Include: []string{"/.*/"},
			Exclude: []string{"/^hello$/"},
		},
	}
)

func TestNew(t *testing.T) {

	var input interface{}
	var err error

	cases := map[string]struct {
		EnvAccountIDSet, EnvAccessTokenSet, ExpectedError bool
	}{
		"No env vars":        {false, false, true},
		"No accessToken env": {true, false, true},
		"No accountID env ":  {false, true, true},
		"Env vars set":       {true, true, false},
	}

	os.Setenv(envAccountID, "username")
	os.Setenv(envAccessToken, "token")
	defer os.Unsetenv(envAccountID)
	defer os.Unsetenv(envAccessToken)

	for tcName, tc := range cases {
		tcSource := source
		if !tc.EnvAccountIDSet {
			tcSource.AccountID = "FAKE_ACCOUNT_ID"
		}
		if !tc.EnvAccessTokenSet {
			tcSource.AccessToken = "FAKE_ACCESS_TOKEN"
		}
		input, err = bbcloud.New(tcSource)

		_, typeOK := input.(plugins.Input)

		actualError := (err != nil)
		errorOK := (actualError == tc.ExpectedError)
		if !errorOK {
			t.Errorf("%v - Expected error: %v | Actual Error: %v", tcName, tc.ExpectedError, actualError)
		}
		if !typeOK {
			t.Errorf("%v - Plugin type check failed", tcName)
		}
	}
}

func TestCredentials(t *testing.T) {

	var input plugins.Input

	cases := map[string]struct {
		EnvAccountIDSet, EnvAccessTokenSet, ExpectedError bool
		ExpectedAccountID, ExpectedAccessToken            string
	}{
		"No env vars":        {false, false, true, "", ""},
		"No accessToken env": {true, false, true, "", ""},
		"No accountID env ":  {false, true, true, "", ""},
		"Env vars set":       {true, true, false, "username", "token"},
	}

	os.Setenv(envAccountID, "username")
	os.Setenv(envAccessToken, "token")
	input, err := bbcloud.New(source)
	if err != nil {
		t.Error("Plugin initiation failed")
	}
	os.Unsetenv(envAccountID)
	os.Unsetenv(envAccessToken)

	for tcName, tc := range cases {
		t.Run(tcName, func(t *testing.T) {
			if tc.EnvAccountIDSet {
				os.Setenv(envAccountID, "username")
				defer os.Unsetenv(envAccountID)
			}
			if tc.EnvAccessTokenSet {
				os.Setenv(envAccessToken, "token")
				defer os.Unsetenv(envAccessToken)
			}

			accountID, accessToken, err := input.Credentials()

			idOK := accountID == tc.ExpectedAccountID
			tokenOK := accessToken == tc.ExpectedAccessToken

			actualError := (err != nil)
			errorOK := (actualError == tc.ExpectedError)
			if !(idOK && tokenOK && errorOK) {
				t.Errorf("%v - Expected error: %v | Actual Error: %v", tcName, tc.ExpectedError, actualError)
			}
			os.Unsetenv(envAccountID)
			os.Unsetenv(envAccessToken)
		})
	}
}

func TestAuthenticate(t *testing.T) {

	var input plugins.Input
	var err error

	cases := map[string]struct {
		AccountID, AccessToken        string
		ExpectedResult, ExpectedError bool
	}{
		"Wrong credentials":   {"usern", "tok", false, true},
		"Wrong account ID":    {"kuserna", "token", false, true},
		"Wrong access token":  {"username", "tochen", false, true},
		"Correct credentials": {"username", "token", true, false},
	}

	for tcName, tc := range cases {
		t.Run(tcName, func(t *testing.T) {
			os.Setenv(envAccountID, tc.AccountID)
			os.Setenv(envAccessToken, tc.AccessToken)
			defer os.Unsetenv(envAccountID)
			defer os.Unsetenv(envAccessToken)

			// New
			input, err = bbcloud.New(source)
			if err != nil {
				t.Error("Plugin initiation failed")
			}

			// Set mock client
			mockInput := input.(*bbcloud.Cloud)
			mockInput.API.Teams = &mock.Teams{AccountID: tc.AccountID, AccessToken: tc.AccessToken}
			mockInput.API.Repositories = &mock.Repositories{AccountID: tc.AccountID, AccessToken: tc.AccessToken}

			// Call Authenticate
			actualResult, err := mockInput.Authenticate()
			actualError := (err != nil)

			// Validate
			resultOK := (actualResult == tc.ExpectedResult)
			errorOK := (actualError == tc.ExpectedError)
			if !(resultOK && errorOK) {
				t.Errorf("%v - Expected error: %v | Actual Error: %v", tcName, tc.ExpectedError, actualError)
			}
		})
	}

	// Test for different Kind values
	os.Setenv(envAccountID, "username")
	os.Setenv(envAccessToken, "token")
	defer os.Unsetenv(envAccountID)
	defer os.Unsetenv(envAccessToken)

	kindCases := map[string]struct {
		Kind                          string
		ExpectedResult, ExpectedError bool
	}{
		"Wrong username":       {"user/unamebad", false, true},
		"Wrong project key":    {"project/NOYA", false, true},
		"Unsupported kind":     {"proj/GIGA", false, true},
		"Correct username":     {"user/username", true, false},
		"Correct project name": {"project/TEST", true, false},
	}

	for tcName, tc := range kindCases {
		t.Run(tcName, func(t *testing.T) {

			tcSource := source
			tcSource.Kind = tc.Kind

			// New
			input, err = bbcloud.New(tcSource)
			if err != nil {
				t.Error("Plugin initiation failed")
			}

			// Set mock client
			mockInput := input.(*bbcloud.Cloud)
			mockInput.API.Teams = &mock.Teams{AccountID: "username", AccessToken: "token"}
			mockInput.API.Repositories = &mock.Repositories{AccountID: "username", AccessToken: "token"}

			// Call Authenticate()
			actualResult, err := mockInput.Authenticate()
			actualError := (err != nil)

			// Validate
			resultOK := (actualResult == tc.ExpectedResult)
			errorOK := (actualError == tc.ExpectedError)
			if !(resultOK && errorOK) {
				t.Error("Authentication test failed")
			}
		})
	}
}

func TestRepositories(t *testing.T) {
	var input plugins.Input
	var err error

	cases := map[string]struct {
		Kind                          string
		ExpectedResult, ExpectedError bool
	}{
		"Wrong username":       {"user/unamebad", false, true},
		"Wrong project key":    {"project/NOYA", false, false},
		"Unsupported kind":     {"proj/GIGA", false, true},
		"Correct username":     {"user/username", true, false},
		"Correct project name": {"project/TEST", true, false},
	}

	os.Setenv(envAccountID, "username")
	os.Setenv(envAccessToken, "token")
	defer os.Unsetenv(envAccountID)
	defer os.Unsetenv(envAccessToken)

	for tcName, tc := range cases {
		t.Run(tcName, func(t *testing.T) {

			tcSource := source
			tcSource.Kind = tc.Kind

			// New
			input, err = bbcloud.New(tcSource)
			if err != nil {
				t.Error("Plugin initiation failed")
			}

			// Set mock client
			mockInput := input.(*bbcloud.Cloud)
			mockInput.API.Teams = &mock.Teams{AccountID: "username", AccessToken: "token"}
			mockInput.API.Repositories = &mock.Repositories{AccountID: "username", AccessToken: "token"}

			// Call Repositories()
			result, err := mockInput.Repositories(true)

			actualResult := (result != nil) && (len(result) != 0)
			actualError := (err != nil)

			// Validate
			resultOK := (actualResult == tc.ExpectedResult)
			errorOK := (actualError == tc.ExpectedError)
			if !(resultOK && errorOK) {
				t.Error("Repositories() test failed")
			}
		})
	}
}
