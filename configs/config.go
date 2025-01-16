package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// See init()
var sharedConfig *Config

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	Client struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	SecretKeeper string `yaml:"secretKeeper"`

	CORS struct {
		AllowedOrigins []string `yaml:"allowedOrigins,flow"`
	}
}

// Fill sharedConfig with config obtained from the file.
// This init will run across files and packages
// that use config package.
func init() {
	if sharedConfig != nil {
		return
	}

	err := error(nil)
	sharedConfig, err = getConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func Configuration() *Config {
	return sharedConfig
}

func getConfig() (*Config, error) {
	configPath, err := parseFlag()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Using config with path:", configPath)

	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	log.Printf("Using these configs: %+v\n", config)
	return config, nil
}

func parseFlag() (string, error) {
	var configPath string

	// Capability to use another config file
	// by supplying -config flag.
	// Default to config.yaml
	flag.StringVar(&configPath, "config", "./config.yaml", "path to config file")

	flag.Parse()

	if err := validateConfigPath(configPath); err != nil {
		return "", err
	}

	return configPath, nil
}

func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
