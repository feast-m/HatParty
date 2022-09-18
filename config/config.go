package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	MongoConnectionString string `yaml:"MongoConnectionString"`
	CleaningTimeInHours   int    `yaml:"CleaningTimeInHours"`
	MaxHatsPerParty       int    `yaml:"MaxHatsPerParty"`
	InitNumberOfHats      int    `yaml:"InitNumberOfHats"`
}

func newConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
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

func parseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := validateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}

func ParseConfig() *Config {
	cfgPath, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := newConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

var Cfg *Config = ParseConfig()
