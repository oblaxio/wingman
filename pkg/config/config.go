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
	Version       float64                  `yaml:"version"`
	Module        string                   `yaml:"module"`
	Env           map[string]string        `yaml:"env,omitempty"`
	BuildDir      string                   `yaml:"build_dir"`
	Watchers      Watchers                 `yaml:"watchers"`
	Services      map[string]ServiceConfig `yaml:"services"`
	Proxy         Proxy                    `yaml:"proxy"`
	ServiceGroups map[string][]string      `yaml:"service_groups"`
}

type Watchers struct {
	IncludeDirs  []string `yaml:"include_dirs"`
	ExcludeDirs  []string `yaml:"exclude_dirs"`
	IncludeFiles []string `yaml:"include_files"`
	ExcludeFiles []string `yaml:"exclude_files"`
}

type Proxy struct {
	Enabled     bool         `yaml:"enabled"`
	Port        int          `yaml:"port"`
	Address     string       `yaml:"address"`
	APIPrefix   string       `yaml:"api_prefix"`
	LogRequests bool         `yaml:"log_requests"`
	Storage     ProxyStorage `yaml:"storage,omitempty"`
	SPA         ProxySPA     `yaml:"spa,omitempty"`
	Static      ProxyStatic  `yaml:"static,omitempty"`

	// PublicAssets ProxyPublicAssets `yaml:"public_assets,omitempty"`
}

type ProxyStorage struct {
	Enabled bool   `yaml:"enabled"`
	Prefix  string `yaml:"prefix"`
	Bucket  string `yaml:"bucket"`
	Service string `yaml:"service"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type ProxySPA struct {
	Enabled bool   `yaml:"enabled"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type ProxyStatic struct {
	Enabled bool   `yaml:"enabled"`
	Dir     string `yaml:"dir"`
	Index   string `yaml:"index"`
}

// type ProxyPublicAssets struct {
// 	Enabled            bool   `yaml:"enabled"`
// 	Path               string `yaml:"path"`
// 	Dir                string `yaml:"dir"`
// 	ServiceNamePattern string `yaml:"service_name_pattern"`
// }

type ServiceConfig struct {
	Entrypoint     string            `yaml:"entrypoint"`
	Executable     string            `yaml:"executable"`
	Env            map[string]string `yaml:"env"`
	ProxyType      string            `yaml:"proxy_type"`
	ProxyHandle    string            `yaml:"proxy_handle"`
	ProxyStaticDir string            `yaml:"proxy_static_dir"`
	ProxyIndex     string            `yaml:"proxy_index"`
	ProxyPort      int               `yaml:"proxy_port"`
	ProxyAddress   string            `yaml:"proxy_address"`
	LDFlags        map[string]string `yaml:"ldflags"`
}

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
