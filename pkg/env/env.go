package env

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Var struct {
	Key   string
	Value string
}

type Vars []*Var

func FromFile(envFile string) (map[string]string, error) {
	out := make(map[string]string)
	// check if file exists
	if len(envFile) == 0 {
		return nil, errors.New("no env file specified")
	}
	if _, err := os.Stat(envFile); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("file does't exist %s: %s", envFile, err.Error())
	}
	// read file contents
	file, err := os.ReadFile(envFile)
	if err != nil {
		return nil, errors.New("import of env file failed " + envFile)
	}
	// parse file contents
	content := string(file)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				out[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}
	return out, nil
}
