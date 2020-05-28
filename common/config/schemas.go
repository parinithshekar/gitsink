package config

type Sync struct {
	Type   string `yaml:"type"`
	Period int    `yaml:"period_seconds"`
}

type Filters struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude,omitempty"`
}

type Repositories struct {
	Filters Filters `yaml:"filters"`
}

type Source struct {
	Type         string       `yaml:"type"`
	BaseURL      string       `yaml:"base_url,omitempty"`
	AccountID    string       `yaml:"account_id"`
	AccessToken  string       `yaml:"access_token"`
	Kind         string       `yaml:"kind"`
	Repositories Repositories `yaml:"repos"`
}

type BranchModifier struct {
	Name string `yaml:"name"`
	Match string `yaml:"match"`
	Prefix string `yaml:"prefix,omitempty"`
	Rename string `yaml:"rename,omitempty"`
}

type Target struct {
	Type         string       `yaml:"type"`
	BaseURL      string       `yaml:"base_url"`
	AccountID    string       `yaml:"account_id"`
	AccessToken  string       `yaml:"access_token"`
	Kind         string       `yaml:"kind"`
	BranchModifiers []BranchModifier `yaml:"branch_modifiers"`
}

type Integration struct {
	Name    string `yaml:"name"`
	Enabled bool   `yaml:"enabled"`
	Sync    Sync   `yaml:"sync"`
	Source  Source `yaml:"source"`
	Target Target `yaml:"target"`
}

type Config struct {
	Integrations []Integration `yaml:"integrations"`
}
