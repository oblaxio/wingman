package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/mod/modfile"
	"gopkg.in/yaml.v3"
)

const DefaultConfigFile = "wingman.yaml"

var configInstance *Config

type Config struct {
	// Configuration file version
	Version float64 `yaml:"version"`
	// Go module name
	Module string `yaml:"module"`
	// A key-value list of environment variables
	Env map[string]string `yaml:"env,omitempty"`
	// List of environment variable (.env) files
	EnvFiles []string `yaml:"env_files,omitempty"`
	// The build output directory
	BuildDir string `yaml:"build_dir"`
	// File and directory watchers
	Watchers Watchers `yaml:"watchers"`
	// Service configurations
	Services map[string]ServiceConfig `yaml:"services"`
	// Proxy server configuration
	Proxy Proxy `yaml:"proxy"`
	// Service groups
	ServiceGroups map[string][]string `yaml:"service_groups"`
}

type Watchers struct {
	// Directories to include/exclude from watching
	IncludeDirs []string `yaml:"include_dirs"`
	// Directories to exclude from watching
	ExcludeDirs []string `yaml:"exclude_dirs"`
	// Files to include/exclude from watching
	IncludeFiles []string `yaml:"include_files"`
	// Files to exclude from watching
	ExcludeFiles []string `yaml:"exclude_files"`
}

type Proxy struct {
	// Tells whether to start the proxy server
	Enabled bool `yaml:"enabled"`
	// The port on which the proxy server will listen
	Port int `yaml:"port"`
	// The address on which the proxy server will listen
	Address string `yaml:"address"`
	// The API prefix to use for routing
	APIPrefix string `yaml:"api_prefix"`
	// Whether to log incoming requests to the terminal
	LogRequests bool `yaml:"log_requests"`
}

type ServiceConfig struct {
	// Relative location of the service
	Entrypoint string `yaml:"entrypoint"`
	// Name of the output executable
	Executable string `yaml:"executable"`
	// A key-value list of environment variables
	Env map[string]string `yaml:"env"`
	// List of environment variable (.env) files
	EnvFiles []string `yaml:"env_files"`
	// The type of the proxy. Can be "service" "storage"
	ProxyType string `yaml:"proxy_type"`
	// The route handle/prefix to proxy requests to this service
	ProxyHandle string `yaml:"proxy_handle"`
	// Works only with "storage" proxy type. The directory to serve static files from
	ProxyStaticDir string `yaml:"proxy_static_dir"`
	// The index file to serve for storage applications
	ProxyIndex string `yaml:"proxy_index"`
	// The port of the service to proxy requests to
	ProxyPort int `yaml:"proxy_port"`
	// The address of the service to proxy requests to
	ProxyAddress string `yaml:"proxy_address"`
	// Rewrite (replace) part of the route when proxying. Format: requested_route:target_route
	ProxyRouteRewrite string `yaml:"proxy_rewrite"`
	// Flags to pass to the go build command
	LDFlags map[string]string `yaml:"ldflags"`
	// List of services this service depends on
	DependsOn []string `yaml:"depends_on"`
}

// Creates a new default config
func NewConfig() *Config {
	f, err := os.ReadFile("go.mod")
	if err != nil {
		panic(err)
	}
	// parse the go.mod file
	m, err := modfile.Parse("go.mod", f, nil)
	if err != nil {
		panic(err)
	}
	// create the default build directory if it doesn't exist
	if _, err := os.Stat("./bin"); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir("./bin", 0755); err != nil {
			panic(err)
		}
	}
	c := &Config{
		Version:  1.0,
		Module:   m.Module.Mod.Path,
		BuildDir: "bin",
		Watchers: Watchers{
			IncludeFiles: []string{"*.go"},
			ExcludeFiles: []string{"test_*.go"},
			IncludeDirs:  []string{"."},
			ExcludeDirs:  []string{"bin", "vendor", "docs"},
		},
		Proxy: Proxy{
			Enabled:     true,
			Port:        8080,
			Address:     "127.0.0.1",
			APIPrefix:   "/api",
			LogRequests: true,
		},
		Services: map[string]ServiceConfig{},
	}
	return c
}

// Creates the config file on disk. Errors if the file already exists
func (c *Config) Create() error {
	if _, err := os.Stat(DefaultConfigFile); errors.Is(err, os.ErrNotExist) {
		fh, err := os.Create(DefaultConfigFile)
		if err != nil {
			panic(err)
		}
		defer fh.Close()
		enc := yaml.NewEncoder(fh)
		enc.SetIndent(2)
		if err := enc.Encode(c); err != nil {
			panic(err)
		}
		return nil
	}
	return errors.New("file already exists")
}

// Reads the config from a file. Errors if the file doesn't exist or is malformed
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

// Path returns the directory path of the config file
func Path(configFilePath string) string {
	pathParts := strings.Split(configFilePath, "/")
	if len(pathParts) > 0 {
		return strings.Join(pathParts[:len(pathParts)-1], "/")
	}
	return ""
}
