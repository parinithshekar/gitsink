package git_test

import (
	"os"
	"testing"

	config "github.com/parinithshekar/gitsink/common/config"
	bbserver "github.com/parinithshekar/gitsink/plugins/input/bitbucket/server"
	git "github.com/parinithshekar/gitsink/plugins/output/git"
	ghpublic "github.com/parinithshekar/gitsink/plugins/output/github/public"
)

var (
	envSourceAccountID   = "TEST_BBSERVER_ACCOUNT_ID"
	envSourceAccessToken = "TEST_BBSERVER_ACCESS_TOKEN"
	envTargetAccountID   = "TEST_GHPUBLIC_ACCOUNT_ID"
	envTargetAccessToken = "TEST_GHPUBLIC_ACCESS_TOKEN"

	source config.Source = config.Source{
		Type:        "bitbucket-server",
		BaseURL:     "https://bitbucket-test.company.com",
		AccountID:   envSourceAccountID,
		AccessToken: envSourceAccessToken,
		Kind:        "user/username",
		Repositories: config.Repositories{
			Include: []string{".*-suffix"},
			Exclude: []string{"prefix-.*"},
		},
	}

	target config.Target = config.Target{
		Type:        "github-public",
		AccountID:   envTargetAccountID,
		AccessToken: envTargetAccessToken,
		Kind:        "user/username",
	}
)

func TestNew(t *testing.T) {
	os.Setenv(envSourceAccountID, "username")
	os.Setenv(envSourceAccessToken, "token")
	os.Setenv(envTargetAccountID, "username")
	os.Setenv(envTargetAccessToken, "token")
	defer os.Unsetenv(envSourceAccountID)
	defer os.Unsetenv(envSourceAccessToken)
	defer os.Unsetenv(envTargetAccountID)
	defer os.Unsetenv(envTargetAccessToken)

	input, _ := bbserver.New(source)
	output, _ := ghpublic.New(target)

	_ = git.New(input, output, "test-integration")
}
