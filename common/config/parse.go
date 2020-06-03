package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Parse Reads the YAML config file and returns a corresponding map
func Parse() Config {
	config := Config{}

	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Printf("YAML file read: %v\n", err)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		fmt.Printf("Unmarshall: %v\n", err)
	}
	return config
}
