/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [target]",
	Short: "Installs the mcp gateway proxy",
	Long: `Installs the MCP gateway proxy in target environment.
	example: \"install claude\" install the MCP Gate as an MCP Server in the local Claude-Desktop on the machine
	`,
}

// claudeCmd represents the claude command
var claudeCmd = &cobra.Command{
	Use:   "claude",
	Short: "local Claude-Desktop",
	Long: `target for this operation is the Local Claude-Desktop.
	This command lookss
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Installing mcp gate into local Claude Desktop")
		userconfigdir, _ := os.UserConfigDir()
		configdir := path.Join(userconfigdir, "Claude")
		configfilename := path.Join(configdir, "claude_desktop_config.json")
		if config, found := readClaudeDesktopConfig(configfilename); found {
			if err := backupClaudeDesktopConfig(configfilename); err != nil {
				log.Fatalf("Error backing up Claude Desktop config file:%s\n", err)
			}
			addMCPGateToClaudeDesktopConfig(config)
			err := saveClaudeDesktopConfig(configfilename, config)
			if err != nil {
				log.Println("error writing config file: %w", err)
			}
		} else {
			log.Printf("Claude desktop config not found in %s. Maybe Claude Desktop not installed ?\n", configfilename)
		}

	},
}

func backupClaudeDesktopConfig(fileName string) error {
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

func saveClaudeDesktopConfig(fileName string, config map[string]interface{}) error {
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

func addMCPGateToClaudeDesktopConfig(config map[string]interface{}) error {

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
			"args":    []string{"server"},
		}
	}
	return nil
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

func readClaudeDesktopConfig(fileName string) (map[string]interface{}, bool) {
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

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.AddCommand(claudeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
