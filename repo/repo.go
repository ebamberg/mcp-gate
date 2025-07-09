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
	Command      string   `yaml:"command"`
	Args         []string `yaml:"args"`
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
