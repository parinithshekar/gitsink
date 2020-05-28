package public

import (

	"os"
	"fmt"
	"strings"
	"context"

	"golang.org/x/oauth2"
	"github.com/google/go-github/v31/github"
	
	"github.com/parinithshekar/github-migration-cli/common/config"
	logger "github.com/parinithshekar/github-migration-cli/wrap/logrus/v1"
)
// Public struct defines fields in github-public object
type Public struct {
	accountID string
	accessToken string
	kind string
	api *github.Client
	ctx context.Context
}

func (public * Public) setAPIClient() {
	ctx := context.Background()

	accessToken := os.Getenv(public.accessToken)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	public.api = github.NewClient(tc)
	public.ctx = ctx
}

// New returns a new github-public object
func New(target config.Target) (*Public) {
	var public *Public = new(Public)

	// Check if env variables mentioned in config file exist
	// check account ID env variable
	_, exists := os.LookupEnv(target.AccountID)
	if !exists {
		logger.New().Errorf("Account ID not found")
		os.Exit(1)
	} else {
		public.accountID = target.AccountID
	}
	// check access token env variable
	_, exists = os.LookupEnv(target.AccessToken)
	if !exists { 
		logger.New().Errorf("Access Token not found")
		os.Exit(1)
	} else {
		public.accessToken = target.AccessToken
	}

	public.kind = target.Kind

	public.setAPIClient()
	return public
}

// Authenticate checks the account ID and access tokens' validity for the kind defined
func (public Public) Authenticate() (bool, error) {
	kindSplit := strings.SplitN(public.kind, "/", 2)
	kindType := kindSplit[0]
	kindKey := kindSplit[1]

	accountID := os.Getenv(public.accountID)

	switch kindType {
	case "org":
		// returns true only if the user is a member of the org, can create repositories
		// returns Membership, Response, error
		_, _, err := public.api.Organizations.GetOrgMembership(public.ctx, accountID, kindKey)
		if err != nil {
			return false, err
		}
		return true, nil
	
	case "user":
		// return true if the authenticated user from env variables is the same user mentioned in config
		user, _, err := public.api.Users.Get(context.Background(), "")
		if err != nil {
			fmt.Println("USER AUTH: ", err.Error())
			return false, err
		} else if kindKey != *user.Login {
			return false, fmt.Errorf("Kind username does not match account ID. Unable to push to this target")
		}
		return true, nil
	
	default:
		// Mentioned kind is unsupported
		return false, fmt.Errorf("Unsupported kind: %v", kindType)
	}
}