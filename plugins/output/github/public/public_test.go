package public_test

import (
	"os"
	"testing"

	config "github.com/parinithshekar/gitsink/common/config"
	plugins "github.com/parinithshekar/gitsink/plugins/interfaces"
	ghpublic "github.com/parinithshekar/gitsink/plugins/output/github/public"
)

var (
	envAccountID   = "TEST_GHPUBLIC_ACCOUNT_ID"
	envAccessToken = "TEST_GHPUBLIC_ACCESS_TOKEN"

	target config.Target = config.Target{
		Type:        "github-public",
		AccountID:   envAccountID,
		AccessToken: envAccessToken,
		Kind:        "user/username",
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
		tcTarget := target
		if !tc.EnvAccountIDSet {
			tcTarget.AccountID = "FAKE_ACCOUNT_ID"
		}
		if !tc.EnvAccessTokenSet {
			tcTarget.AccessToken = "FAKE_ACCESS_TOKEN"
		}
		input, err = ghpublic.New(tcTarget)

		_, typeOK := input.(plugins.Output)

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

	var input plugins.Output

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
	input, err := ghpublic.New(target)
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
