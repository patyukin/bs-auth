package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func LoadYaml() (*Config, error) {
	appEnv := os.Getenv(EnvFilePath)

	var configFilePath string
	// проверка возможных вариантов для appEnv
	switch appEnv {
	case "local-dev":
		configFilePath = "config.local-dev.yaml"
	case "local":
		configFilePath = "config.local.yaml"
	case "test":
		configFilePath = "config.test.yaml"
	case "test-dev":
		configFilePath = "config.test-dev.yaml"
	case "prod":
		configFilePath = "config.yaml"
	default:
		panic("config not found")
	}

	cfg := Config{}

	// Open config file
	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(file)

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err = d.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
