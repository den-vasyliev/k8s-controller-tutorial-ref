package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/yourusername/k8s-controller-tutorial/pkg/api"
)

// NewMCPServer creates and configures a new MCP server for FrontendPage tools
func NewMCPServer(serverName, version string) *server.MCPServer {
	s := server.NewMCPServer(
		serverName,
		version,
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// List tool
	listTool := mcp.NewTool("list_frontendpages",
		mcp.WithDescription("List all FrontendPage resources"),
	)
	// Create tool
	createTool := mcp.NewTool("create_frontendpage",
		mcp.WithDescription("Create a new FrontendPage resource"),
		mcp.WithString("name", mcp.Description("Name of the FrontendPage")),
		mcp.WithString("contents", mcp.Description("HTML contents")),
		mcp.WithString("image", mcp.Description("Container image")),
		mcp.WithNumber("replicas", mcp.Description("Number of replicas")),
	)
	// TODO: Add update and delete tools as needed

	s.AddTool(listTool, listFrontendPagesHandler)
	s.AddTool(createTool, createFrontendPageHandler)
	// TODO: Register update/delete handlers

	return s
}

func listFrontendPagesHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if api.FrontendAPI == nil {
		return mcp.NewToolResultText("FrontendAPI is not initialized"), nil
	}
	docs, err := api.FrontendAPI.ListFrontendPagesRaw(ctx)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error listing FrontendPages: %v", err)), nil
	}
	jsonBytes, err := json.MarshalIndent(docs, "", "  ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error marshaling result: %v", err)), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func createFrontendPageHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// TODO: Integrate with your API logic to create a FrontendPage
	return mcp.NewToolResultText("Created FrontendPage (stub)"), nil
}
