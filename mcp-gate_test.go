package main

import (
	"log"
	"testing"

	"github.com/spf13/viper"
)

func TestReadConfig(t *testing.T) {
	initConfig()

	if viper.GetString("app.name") != "mcp-gate" {
		t.Errorf("Expected app name 'mcp-gate', got '%s'", viper.GetString("app.name"))
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
