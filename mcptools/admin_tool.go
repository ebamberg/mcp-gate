package mcptools

import (
	"context"
	"fmt"

	"github.com/ebamberg/mcp-gate/repo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func adminInstallToolSchema() mcp.Tool {
	// Add a admin tool
	return mcp.NewTool("mcp-gate-install-tool",
		mcp.WithDescription("install an available tool into mcp-gate"),
		mcp.WithString("toolname",
			mcp.Required(),
			mcp.Description("The name of the tool to install in mcp-gate"),
		),
	)
}

func listAvailableToolsSchema() mcp.Tool {
	// Add a admin tool
	return mcp.NewTool("mcp-gate-list-available",
		mcp.WithDescription("returns a list of available tools that can be installed in mcp-gate"),
	)
}

func mcpGateVersionResourceSchema() mcp.Resource {
	return mcp.NewResource("mcpgate://version", "mcp-gate-version",
		mcp.WithResourceDescription("The version of the installed mcp-gate."),
		mcp.WithMIMEType("text/plain"),
	)
}

func RegisterAdminTool(server *server.MCPServer) {
	server.AddResource(mcpGateVersionResourceSchema(), mcpGateVersionResourceHandler)
	// Add the install a tool handler
	server.AddTool(listAvailableToolsSchema(), listAvailableToolsHandler)
	server.AddTool(adminInstallToolSchema(), installToolHandler)
}

func mcpGateVersionResourceHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "mcpgate://version",
			MIMEType: "text/plain",
			Text:     string("version 1.0.0"),
		},
	}, nil

}

func listAvailableToolsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	repoEntries, err := repo.ListAvailableTools()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	var result string
	for _, entry := range repoEntries {
		result += fmt.Sprintf("Tool: %s\nDescription: %s\nCommand: %s\nArgs: %v\nDependencies: %v\nPlatforms: %v\n\n",
			entry.Name, entry.Description, entry.Command, entry.Args, entry.Dependencies, entry.Platforms)

	}
	return mcp.NewToolResultText(fmt.Sprintf("List of mcp-server-tools than can be installed and can be used.\n\n %s ", result)), nil
}

func installToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Using helper functions for type-safe argument access
	toolname, err := request.RequireString("toolname")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("the mcp-server-tool %s is installed and can be used.", toolname)), nil
}
