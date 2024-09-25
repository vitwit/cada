package chainclient

import (
	"fmt"
	"os"
)

// Config represents the structure of the config.yaml file
type Config struct {
	ValidatorKey   string `yaml:"validatorKey"`   // It is the keyname of the cosmos validator account to sign the transactions
	KeyringBackend string `yaml:"keyringBackend"` // It is the type of the keyring to sign the transactions related to submission of data to avail
}

func GetClientConfig() Config {
	config := Config{ //TODO: think about better approach
		ValidatorKey:   os.Getenv("VALIDATOR_KEY"),
		KeyringBackend: os.Getenv("KEYRING_BACKEND"),
	}

	// Output the configuration values
	fmt.Printf("ValidatorKey: %s\n", config.ValidatorKey)
	fmt.Printf("KeyringType: %s\n", config.KeyringBackend)
	return config
}
