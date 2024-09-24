package chainclient

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of the config.yaml file
type Config struct {
	ValidatorKey   string `yaml:"validatorKey"`   // It is the keyname of the cosmos validator account to sign the transactions
	KeyringBackend string `yaml:"keyringBackend"` // It is the type of the keyring to sign the transactions related to submission of data to avail
}

func GetClientConfig() Config {
	// Read the YAML configuration file
	wd, err := os.Getwd()
	fmt.Println("working dirrr........", wd)
	configpath := fmt.Sprintf(wd + "/chainclient/config.yaml")
	fmt.Println("comnfig pathhhhh....", configpath)
	data, err := os.ReadFile(configpath)
	if err != nil {
		log.Println("Error reading YAML file: %v", err)
	}

	// Initialize a Config struct
	var config Config

	// Unmarshal the YAML data into the Config struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Println("Error unmarshalling YAML: %v", err)
	}

	// Output the configuration values
	fmt.Printf("ValidatorKey: %s\n", config.ValidatorKey)
	fmt.Printf("KeyringType: %s\n", config.KeyringBackend)
	return config
}
