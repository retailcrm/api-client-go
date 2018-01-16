package retailcrm

import (
	"encoding/json"
	"os"
	"fmt"
)

type Configuration struct {
	Url string `json:"url"`
	Key string `json:"key"`
	Ver string `json:"version"`
}

func buildConfiguration() *Configuration {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	return &Configuration {
		configuration.Url,
		configuration.Key,
		configuration.Ver,
	}
}

func UnversionedClient() *Client {
	configuration := buildConfiguration()
	return New(configuration.Url, configuration.Key, "")
}

func VersionedClient() *Client {
	configuration := buildConfiguration()
	return New(configuration.Url, configuration.Key, configuration.Ver)
}

