package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/ebamberg/mcp-gate/repo"
	mcpclient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

type ClientStatus int64

const (
	UNINITIALIZED ClientStatus = iota
	CONNECTED
	STOPPED
	FAILED
)

type Client struct {
	Name           string `json:"name"`
	Status         ClientStatus
	proxied_client *mcpclient.Client
	serverInfo     *mcp.InitializeResult
}

func RegisterMCPTool(config repo.RepositoryEntry) error {

	client, err := NewClient(config)
	if err != nil {
		return fmt.Errorf("Failed to build client for tool %s: %v", config.Name, err)
	}
	err = client.Connect()
	//	schema := buildToolSchema(config)
	return nil
}

func initNotificationHandler(c *mcpclient.Client) {
	// Set up notification handler
	c.OnNotification(func(notification mcp.JSONRPCNotification) {
		log.Printf("Received notification: %s\n", notification.Method)
	})
}

func (client *Client) Connect() error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize the client
	log.Println("Initializing client...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "MCP-Gate",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	client.serverInfo, err = client.proxied_client.Initialize(ctx, initRequest)
	if err != nil {
		client.Status = FAILED
		return fmt.Errorf("Failed to initialize: %v", err)
	}

	// Display server information
	log.Printf("Connected to server: %s (version %s)\n",
		client.serverInfo.ServerInfo.Name,
		client.serverInfo.ServerInfo.Version)
	log.Printf("Server capabilities: %+v\n", client.serverInfo.Capabilities)

	client.Status = CONNECTED

	// List available tools if the server supports them
	if client.serverInfo.Capabilities.Tools != nil {
		fmt.Println("Fetching available tools...")
		toolsRequest := mcp.ListToolsRequest{}
		toolsResult, err := client.proxied_client.ListTools(ctx, toolsRequest)
		if err != nil {
			log.Printf("Failed to list tools: %v", err)
		} else {
			fmt.Printf("Server has %d tools available\n", len(toolsResult.Tools))
			for i, tool := range toolsResult.Tools {
				log.Printf("  %d. %s - %s\n", i+1, tool.Name, tool.Description)
			}
		}
	}

	// List available resources if the server supports them
	if client.serverInfo.Capabilities.Resources != nil {
		fmt.Println("Fetching available resources...")
		resourcesRequest := mcp.ListResourcesRequest{}
		resourcesResult, err := client.proxied_client.ListResources(ctx, resourcesRequest)
		if err != nil {
			log.Printf("Failed to list resources: %v", err)
		} else {
			log.Printf("Server has %d resources available\n", len(resourcesResult.Resources))
			for i, resource := range resourcesResult.Resources {
				log.Printf("  %d. %s - %s\n", i+1, resource.URI, resource.Name)
			}
		}
	}

	log.Println("Client initialized successfully...")
	return nil
}

func (client *Client) Stop() error {
	if client.proxied_client == nil || client.Status == UNINITIALIZED || client.Status == FAILED {
		return fmt.Errorf("Client is not initialized")
	}

	log.Println("Stopping client...")
	if err := client.proxied_client.Close(); err != nil {
		client.Status = FAILED
		return fmt.Errorf("Failed to stop client: %v", err)
	}

	client.Status = STOPPED
	log.Println("Client stopped successfully")
	return nil
}

func NewClient(config repo.RepositoryEntry) (*Client, error) {
	var client *Client
	var err error = nil
	if config.Transport == "ipc" {
		client, err = NewIPCClient(config)
	} else if config.Transport == "http" {
		client, err = NewHTTPStreamingClient(config)
	} else {
		return nil, fmt.Errorf("Unsupported transport type: %s", config.Transport)
	}
	return client, err
}

func NewIPCClient(config repo.RepositoryEntry) (*Client, error) {
	client := &Client{
		Name:           config.Name,
		Status:         UNINITIALIZED,
		proxied_client: nil,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var err error

	log.Println("Initializing stdio ipc client...")

	// Create stdio transport with verbose logging
	stdioTransport := transport.NewStdio(config.Command, nil, config.Args...)

	// Create client with the transport
	client.proxied_client = mcpclient.NewClient(stdioTransport)

	// Start the client
	if err = client.proxied_client.Start(ctx); err != nil {
		client.Status = FAILED
		return client, fmt.Errorf("Failed to start mcp client: %v", err)
	}

	// Set up logging for stderr if available
	if stderr, ok := mcpclient.GetStderr(client.proxied_client); ok {
		go func() {
			buf := make([]byte, 4096)
			for {
				n, err := stderr.Read(buf)
				if err != nil {
					if err != io.EOF {
						log.Printf("Error reading stderr: %v", err)
					}
					return
				}
				if n > 0 {
					fmt.Fprintf(os.Stderr, "[Server] %s", buf[:n])
				}
			}
		}()
	} else {
		log.Println("%s: No stderr available for logging", config.Name)
	}

	return client, nil
}

func NewHTTPStreamingClient(config repo.RepositoryEntry) (*Client, error) {
	log.Println("Initializing HTTP client...")

	client := &Client{
		Name:           config.Name,
		Status:         UNINITIALIZED,
		proxied_client: nil,
	}
	// Create HTTP transport
	httpTransport, err := transport.NewStreamableHTTP(*config.URL)
	// NOTE: the default streamableHTTP transport is not 100% identical to the stdio client.
	// By default, it could not receive global notifications (e.g. toolListChanged).
	// You need to enable the `WithContinuousListening()` option to establish a long-live connection,
	// and receive the notifications any time the server sends them.
	//
	//   httpTransport, err := transport.NewStreamableHTTP(*httpURL, transport.WithContinuousListening())
	if err != nil {
		client.Status = FAILED
		return client, fmt.Errorf("Failed to create HTTP transport: %v", err)
	}

	// Create client with the transport
	client.proxied_client = mcpclient.NewClient(httpTransport)
	return client, nil
}

func buildToolSchema(config repo.RepositoryEntry) mcp.Tool {
	// Add a admin tool
	options := []mcp.ToolOption{
		mcp.WithDescription(config.Description),
	}
	schema := mcp.NewTool(config.Name,
		options...,
	)
	return schema
}
