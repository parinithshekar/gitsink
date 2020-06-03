package config

// Sync defines the type and period of the auto sync in config
type Sync struct {
	Type   string `yaml:"type"`
	Period int    `yaml:"period_seconds"`
}

// Filters are regexes or strings to include or exclude repositories
// type Filters struct {
// 	Include []string `yaml:"include"`
// 	Exclude []string `yaml:"exclude,omitempty"`
// }

// Repositories are regexes or strings to include or exclude repositories
type Repositories struct {
	// Filters Filters `yaml:"filters"`
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude,omitempty"`
}

// Source has the fields that describe a source for the sync
type Source struct {
	Type         string       `yaml:"type"`
	BaseURL      string       `yaml:"base_url,omitempty"`
	AccountID    string       `yaml:"account_id"`
	AccessToken  string       `yaml:"access_token"`
	Kind         string       `yaml:"kind"`
	Repositories Repositories `yaml:"repos"`
}

// BranchModifier gives options to modify the branch upon sync
type BranchModifier struct {
	Name   string `yaml:"name"`
	Match  string `yaml:"match"`
	Prefix string `yaml:"prefix,omitempty"`
	Rename string `yaml:"rename,omitempty"`
}

// Target has teh fields that describe a target for the sync
type Target struct {
	Type            string           `yaml:"type"`
	BaseURL         string           `yaml:"base_url"`
	AccountID       string           `yaml:"account_id"`
	AccessToken     string           `yaml:"access_token"`
	Kind            string           `yaml:"kind"`
	BranchModifiers []BranchModifier `yaml:"branch_modifiers"`
}

// Integration defines one integration with all information for sync
type Integration struct {
	Name    string `yaml:"name"`
	Enabled bool   `yaml:"enabled"`
	Sync    Sync   `yaml:"sync"`
	Source  Source `yaml:"source"`
	Target  Target `yaml:"target"`
}

// Config is the parent that defines the config file format
type Config struct {
	Integrations []Integration `yaml:"integrations"`
}
