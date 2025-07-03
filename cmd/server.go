/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/ebamberg/mcp-gate/server"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start the MCP server",
	Long:  `start the MCP Gate proxy as a server and allows Client to connect.`,
	Run: func(cmd *cobra.Command, args []string) {
		redirectToStderr, _ := cmd.Flags().GetBool("redirect-to-stderr")
		if redirectToStderr {
			redirectLoggingToStdErr()
		} else {
			redirectLoggingToFile()
		}
		log.Println("Start MCP Gate server")
		server.StartServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().BoolP("redirect-to-stderr", "", false, "whether to redirect alll log output to stderr. This is useful when the tool runs locally in Claude Desktop to redirct logging to the client log folder.")
}

func redirectLoggingToFile() {
	// Redirect log output to a file

	f, err := os.OpenFile("mcp_gate.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
}

func redirectLoggingToStdErr() {
	// Redirect log output to stderr
	log.SetOutput(os.Stderr)
}
