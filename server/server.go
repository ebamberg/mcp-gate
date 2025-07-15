package server

import (
	"log"

	"github.com/mark3labs/mcp-go/server"
)

func NewServer() *server.MCPServer {
	s := server.NewMCPServer(
		"MCP Gate",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	return s
}

func StartServer(s *server.MCPServer) {
	// Start the server
	log.Println("listener on stdin/stdout")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}
