/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/ebamberg/mcp-gate/integration"
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

		if config, found := integration.ReadClaudeDesktopConfig(configfilename); found {
			if err := integration.BackupClaudeDesktopConfig(configfilename); err != nil {
				log.Fatalf("Error backing up Claude Desktop config file:%s\n", err)
			}
			integration.AddMCPGateToClaudeDesktopConfig(config)
			err := integration.SaveClaudeDesktopConfig(configfilename, config)
			if err != nil {
				log.Println("error writing config file: %w", err)
			}
		} else {
			log.Printf("Claude desktop config not found in %s. Maybe Claude Desktop not installed ?\n", configfilename)
		}

	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.AddCommand(claudeCmd)
}
