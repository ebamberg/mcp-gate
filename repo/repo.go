package repo

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed repo_tools.yaml
var repo_tools_yaml []byte

type RepositoryEntry struct {
	Name         string   `yaml:"name"`
	Description  string   `yaml:"description"`
	Transport    string   `yaml:"transport"`
	URL          *string  `yaml:"url,omitempty"`     // Optional, used for HTTP transport
	Command      string   `yaml:"command,omitempty"` // Optional, used for ipc transport
	Args         []string `yaml:"args,omitempty"`    // Optional, used for ipc transport
	Dependencies []string `yaml:"dependencies,omitempty"`
	Platforms    []string `yaml:"platforms,omitempty"`
}

func ListAvailableTools() ([]RepositoryEntry, error) {
	return loadEmbeddedRepo()
}

func loadEmbeddedRepo() ([]RepositoryEntry, error) {
	var entries []RepositoryEntry
	err := yaml.Unmarshal(repo_tools_yaml, &entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
