package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type Config struct {
	SwimlaneRegion      string
	SwimlaneAccountId   string
	SwimlaneAccessToken string
}

// FQDN returns the fully qualified domain name for the Swimlane region.
func (c *Config) FQDN() string {
	return fmt.Sprintf("%s.swimlane.app", c.SwimlaneRegion)
}

// GetConfigDir returns the directory where SwimPeek configuration files are stored.
func GetConfigDir(createDir bool) (string, error) {
	// Check if the SWIMPEEK_CONFIG_DIR environment variable is set
	// If not, use the default cfgDir in the user's home directory
	cfgDir, exists := os.LookupEnv("SWIMPEEK_CONFIG_DIR")
	if !exists {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		cfgDir = path.Join(homeDir, ".swimpeek")
	}

	// Check if the directory exists
	if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
		// If the directory does not exist and createDir is true, create it
		if createDir {
			err := os.MkdirAll(cfgDir, 0750)
			if err != nil {
				return "", fmt.Errorf("failed to create config directory: %s: %w", cfgDir, err)
			}
		} else {
			// If the directory does not exist and createDir is false, return an error
			return "", fmt.Errorf("config directory does not exist: %s: %w", cfgDir, os.ErrNotExist)
		}
	} else if err != nil {
		// If there was an error checking the directory, return it
		return "", fmt.Errorf("failed to stat config directory: %w", err)
	}

	return cfgDir, nil
}

// ReadConfig reads the configuration from the specified directory.
func ReadConfig(cfgDir string) (*Config, error) {
	cfgPath := path.Join(cfgDir, "config.json")

	// Check if the configuration file exists
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		// If the configuration file does not exist, return an error
		return nil, fmt.Errorf("configuration file does not exist: %s: %w", cfgPath, err)
	} else if err != nil {
		// If there was an error checking the file, return it
		return nil, fmt.Errorf("failed to stat config file: %s: %w", cfgPath, err)
	}

	cfg := &Config{}
	file, err := os.Open(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %s: %w", cfgPath, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %s: %w", cfgPath, err)
	}
	return cfg, nil

}

// SaveConfig saves the configuration to a JSON file.
func SaveConfig(cfgDir string, cfg *Config) error {
	cfgPath := path.Join(cfgDir, "config.json")

	// Create or open the configuration file
	file, err := os.Create(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %s: %w", cfgPath, err)
	}
	defer file.Close()

	// Encode the configuration to JSON and write it to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print with indentation
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode config file: %s: %w", cfgPath, err)
	}

	return nil
}
