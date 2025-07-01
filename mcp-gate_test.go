package main

import (
	"log"
	"testing"
)

func TestReadConfig(t *testing.T) {
	config, err := readConfig()
	if err != nil {
		t.Fatalf("Error reading config: %v", err)
	}

	if config.App.Name != "mcp-gate" {
		t.Errorf("Expected app name 'mcp-gate', got '%s'", config.App.Name)
	}

	if config.Namespace != "" {
		t.Errorf("Expected namespace to be empty, got '%s'", config.Namespace)
	}
}

func TestConfigLogging(t *testing.T) {
	configLogging()
	logFlags := log.Flags()
	if logFlags&log.Ldate == 0 || logFlags&log.Ltime == 0 || logFlags&log.Lshortfile == 0 {
		t.Error("Logging flags are not set correctly")
	}

	prefix := log.Prefix()
	if prefix != "MCP-GATE: " {
		t.Errorf("Expected log prefix 'MCP-GATE: ', got '%s'", prefix)
	}
}
