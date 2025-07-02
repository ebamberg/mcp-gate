/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

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
		var configfilename = ""
		switch runtime.GOOS {
		case "windows":
			configfilename = "%APPDATA%\\Claude\\claude_desktop_config.json"
		case "darwin":
			userhome, _ := os.UserHomeDir()
			configfilename = path.Join(userhome, "Library/Application Support/Claude/claude_desktop_config.json")
		default:
			fmt.Println("Operation system not supported")
			os.Exit(-1)
		}
		if config, found := readClaudeDesktopConfig(configfilename); found {
			addMCPGateToClaudeDesktopConfig(config)
		} else {
			log.Printf("Claude desktop config not found in %s. Maybe Claude Desktop not installed ?\n", configfilename)
		}

	},
}

func addMCPGateToClaudeDesktopConfig(config map[string]interface{}) error {

	mcpnode := config["mcpServers"]
	var mcpServers map[string]interface{}
	if mcpnode != nil {
		mcpServers = mcpnode.(map[string]interface{})
	} //else {
	//	mcpServers = map[string]interface{}
	// }
	fmt.Printf("%s\n", mcpnode)
	for _, n := range mcpServers {
		fmt.Println(n)
	}
	return nil
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
