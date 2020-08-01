package cloud_test

import (
	"fmt"
	"os"
	"testing"

	config "github.com/parinithshekar/gitsink/common/config"
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
			Include: []string{".*-suffix"},
			Exclude: []string{"prefix-.*"},
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

	// var input interface{}
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
		fmt.Println(tcName, tc.ExpectedResult)

		os.Setenv(envAccountID, tc.AccountID)
		os.Setenv(envAccessToken, tc.AccessToken)
		defer os.Unsetenv(envAccountID)
		defer os.Unsetenv(envAccessToken)

		// New
		_, err = bbcloud.New(source)
		if err != nil {
			t.Error("Plugin initiation failed")
		}

		// Set mock client
		// Authenticate
		// Validate
	}
}
