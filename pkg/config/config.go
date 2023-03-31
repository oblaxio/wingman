package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const DefaultConfigFile = "wingman.yaml"

var configInstance *Config

type Config struct {
	Version        float64                  `yaml:"version"`
	Module         string                   `yaml:"module"`
	Env            map[string]string        `yaml:"env"`
	BuildDir       string                   `yaml:"build_dir"`
	WatchDir       []string                 `yaml:"watch_dir"`
	DontWatchDir   []string                 `yaml:"dont_watch_dir"`
	WatchFiles     []string                 `yaml:"watch_files"`
	DontWatchFiles []string                 `yaml:"dont_watch_files"`
	Services       map[string]ServiceConfig `yaml:"services"`
}

type ServiceConfig struct {
	Entrypoint string            `yaml:"entrypoint"`
	Executable string            `yaml:"executable"`
	Env        map[string]string `yaml:"env"`
}

func NewConfigFromFile(filePath string) error {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("bad config file '%s', or no config file specified", filePath)
	}
	fh, err := os.ReadFile(filePath)
	if err != nil {
		return errors.New("could not read config file")
	}
	if err := yaml.Unmarshal(fh, &configInstance); err != nil {
		return errors.New("could not parse config file")
	}
	return nil
}

func Get() *Config {
	return configInstance
}

func Path(configFilePath string) string {
	pathParts := strings.Split(configFilePath, "/")
	if len(pathParts) > 0 {
		return strings.Join(pathParts[:len(pathParts)-1], "/")
	}
	return ""
}
