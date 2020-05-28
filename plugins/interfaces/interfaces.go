package interfaces

import "github.com/parinithshekar/github-migration-cli/common"

// Input lists the methods that an input plugin must implement
type Input interface {
	Authenticate() (bool, error)
	Repositories(bool) ([]common.Repository, error)
}

// Output lists the methods that on output plugin must implement
type Output interface {
	Authenticate() (bool, error)
}