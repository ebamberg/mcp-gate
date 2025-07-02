package server

import (
	"log"

	"github.com/mark3labs/mcp-go/server"
)

func StartServer() *server.MCPServer {
	s := server.NewMCPServer(
		"MCP Gate",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Start the server
	log.Println("listener on stdin/stdout")
	if err := server.ServeStdio(s); err != nil {
		log.Printf("Server error: %v\n", err)
	}
	return s
}
