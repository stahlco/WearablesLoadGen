package ssh

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Hosts map[string]Host `yaml:"hosts"`
}

type Host struct {
	IP          string `yaml:"ip,omitempty"`
	Port        int    `yaml:"port,omitempty"`
	Username    string `yaml:"username,omitempty"`
	KeyFile     string `yaml:"key_file,omitempty"`
	KeyPassword string `yaml:"key_password,omitempty"`
}

func ParseYAML(data []byte) (*Config, error) {
	var config Config

	if err := yaml.UnmarshalWithOptions(data, &config, yaml.Strict()); err != nil {
		return nil, err
	}
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(cfg *Config) error {
	if len(cfg.Hosts) == 0 {
		return errors.New("no hosts defined in configuration")
	}

	// Always allow "local" as a valid host alias
	if _, ok := cfg.Hosts["local"]; !ok {
		cfg.Hosts["local"] = Host{}
	}

	for alias, host := range cfg.Hosts {
		if alias == "local" {
			continue
		}
		if host.IP == "" {
			return fmt.Errorf("missing or empty required field: host")
		}
		if host.Username == "" {
			return fmt.Errorf("missing or empty required field: username")
		}
		if host.KeyFile == "" {
			return fmt.Errorf("missing or empty required field: key_file")
		}
		if host.Port == 0 {
			return fmt.Errorf("required field has wrong value: port, has value %d", host.Port)
		}

		if len(host.KeyFile) > 0 && host.KeyFile[0] == '~' {
			if home, err := os.UserHomeDir(); err == nil {
				host.KeyFile = filepath.Join(home, host.KeyFile[1:])
				cfg.Hosts[alias] = host
			}
		}
	}

	return nil
}
