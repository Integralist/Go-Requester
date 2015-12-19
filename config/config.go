package config

import (
	"fmt"
	"io/ioutil"
)

// Parse function reads and returns the contents of provided config file
func Parse(configPath string) []byte {
	config, err := ioutil.ReadFile(configPath)

	if err != nil {
		fmt.Printf("Problem reading the page config: %s\n", err)
		panic(err)
	}

	return config
}
