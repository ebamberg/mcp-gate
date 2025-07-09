package repo

import (
	"testing"
)

func TestListAvailableTools(t *testing.T) {
	tools, err := ListAvailableTools()
	if err != nil {
		t.Fatalf("Failed to list available tools: %v", err)
	}

	if len(tools) == 0 {
		t.Fatal("Expected at least one tool in the repository, but found none.")
	}

	for _, tool := range tools {
		if tool.Name == "" {
			t.Error("Found a tool with an empty name")
		}
		if tool.Command == "" {
			t.Error("Found a tool with an empty command")
		}
	}
}
func TestLoadEmbeddedRepo(t *testing.T) {
	tools, err := loadEmbeddedRepo()
	if err != nil {
		t.Fatalf("Failed to load embedded repository: %v", err)
	}

	if len(tools) == 0 {
		t.Fatal("Expected at least one tool in the embedded repository, but found none.")
	}

	for _, tool := range tools {
		if tool.Name == "" {
			t.Error("Found a tool with an empty name")
		}
		if tool.Command == "" {
			t.Error("Found a tool with an empty command")
		}
	}
}
