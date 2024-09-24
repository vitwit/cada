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
	// Read the YAML configuration file
	// wd, _ := os.Getwd()
	// fmt.Println("working dirrr........", wd)
	// configpath := fmt.Sprintf(wd + "/chainclient/config.yaml")

	// path, err := filepath.Abs("./config.yaml")
	// fmt.Println("path an derrorr.....", path, err)

	// // configPath := config.yaml
	// // fmt.Println("comnfig pathhhhh....", configpath)
	// data, err := os.ReadFile(path)
	// if err != nil {
	// 	log.Println("Error reading YAML file: %v", err)
	// }

	// Initialize a Config struct
	// var config Config

	config := Config{ //TODO: think about better approach
		ValidatorKey:   os.Getenv("VALIDATOR_KEY"),
		KeyringBackend: os.Getenv("KEYRING_BACKEND"),
	}

	// Output the configuration values
	fmt.Printf("ValidatorKey: %s\n", config.ValidatorKey)
	fmt.Printf("KeyringType: %s\n", config.KeyringBackend)
	return config
}
