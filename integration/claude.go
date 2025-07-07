package integration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func BackupClaudeDesktopConfig(fileName string) error {
	backupFileName := fileName + ".bak"
	_, err := os.Stat(backupFileName)
	if err == nil {
		log.Println("Backup file already exists, skipping backup")
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error checking backup file: %w", err)
	}

	fin, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fin.Close()

	fout, err := os.Create(backupFileName)
	if err != nil {
		return fmt.Errorf("error creating backup file: %w", err)
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)
	if err != nil {
		return fmt.Errorf("error creating backup file: %w", err)
	}
	log.Println("Claude Desktop config backed up to", backupFileName)
	return nil
}

func SaveClaudeDesktopConfig(fileName string, config map[string]interface{}) error {
	jsonConfig, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling config: %w", err)
	}

	err = os.WriteFile(fileName, jsonConfig, 0644)
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}
	log.Println("Claude Desktop config saved")
	return nil
}

func AddMCPGateToClaudeDesktopConfig(config map[string]interface{}) error {

	mcpnode := config["mcpServers"]
	var mcpServers map[string]interface{}
	if mcpnode != nil {
		mcpServers = mcpnode.(map[string]interface{})
		// we override the mcp-gate if it already exists for now
		//		for key, _ := range mcpServers {
		//			if key == "mcp-gate" {
		//				log.Println("MCP Gate already installed in Claude Desktop config")
		//				return nil
		//			}
		//		}
		executable, err := getExecutableFilePath()
		if err != nil {
			log.Fatalln("Error getting executable file path:", err)
		}
		mcpServers["mcp-gate"] = map[string]any{
			"command": executable,
			"args":    []string{"server", "--redirect-to-stderr"},
		}
	}
	return nil
}

func ReadClaudeDesktopConfig(fileName string) (map[string]interface{}, bool) {
	datas := map[string]interface{}{}

	file, err := os.ReadFile(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return datas, false
		} else {
			log.Println("unable to read Claude Desktop confilg file", err)
		}
	}

	err = json.Unmarshal(file, &datas)
	if err != nil {
		log.Println("error reading Claude Desktop config file: ", err)
	}
	return datas, true
}

func getExecutableFilePath() (string, error) {
	// get the path to the current executable
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	return execPath, err
}
