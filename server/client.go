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

func (client *Client) addNotificationHandler() error {
	if exit, reason := client.exitOnNotConnected(); exit {
		return reason
	}
	// Set up notification handler
	client.proxied_client.OnNotification(func(notification mcp.JSONRPCNotification) {
		log.Printf("Received notification: %s\n", notification.Method)
	})
	return nil
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

	log.Println("Client initialized successfully...")
	return nil
}

func (client *Client) isConnected() bool {
	return client.proxied_client != nil && client.Status == CONNECTED
}

func (client *Client) exitOnNotConnected() (bool, error) {
	if !client.isConnected() {
		return true, fmt.Errorf("Client is not initialized")
	} else {
		return false, nil
	}
}

func (client *Client) Stop() error {
	if !client.isConnected() {
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

func (client *Client) ListTools() ([]mcp.Tool, error) {

	if exit, reason := client.exitOnNotConnected(); exit {
		return nil, reason
	}

	var tools []mcp.Tool
	var err error = nil

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// List available tools if the server supports them
	if client.serverInfo.Capabilities.Tools != nil {
		log.Println("Fetching available tools...")
		toolsRequest := mcp.ListToolsRequest{}
		toolsResult, err := client.proxied_client.ListTools(ctx, toolsRequest)
		if err != nil {
			log.Printf("Failed to list tools: %v", err)
		} else {
			for _, tool := range toolsResult.Tools {
				tools = append(tools, tool)
			}
		}
	}
	return tools, err
}

func (client *Client) ListResources() ([]mcp.Resource, error) {

	if exit, reason := client.exitOnNotConnected(); exit {
		return nil, reason
	}

	var resources []mcp.Resource
	var err error = nil
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// List available resources if the server supports them
	if client.serverInfo.Capabilities.Resources != nil {
		log.Println("Fetching available resources...")
		resourcesRequest := mcp.ListResourcesRequest{}
		resourcesResult, err := client.proxied_client.ListResources(ctx, resourcesRequest)
		if err != nil {
			log.Printf("Failed to list resources: %v", err)
		} else {
			for _, resource := range resourcesResult.Resources {
				resources = append(resources, resource)
			}
		}
	}
	return resources, err
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
		log.Printf("%s: No stderr available for logging", config.Name)
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
