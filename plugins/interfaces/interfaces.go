package interfaces

import (
	common "github.com/parinithshekar/gitsink/common"
)

// Input lists the methods that an input plugin must implement
type Input interface {
	Authenticate() (bool, error)
	Repositories(bool) ([]common.Repository, error)
	Credentials() (string, string)
}

// Output lists the methods that on output plugin must implement
type Output interface {
	Authenticate() (bool, error)
	SyncCheck([]common.Repository) []common.Repository
	Credentials() (string, string)
}
