package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sashabaranov/go-openai"
	"os"
	"strings"
)

var (
	mcpServerInfo mcp.Implementation
	mcpClient     *client.Client
	openaiClient  *openai.Client
	//go:embed prompt.txt
	systemPromptTemplate string
	messages             []openai.ChatCompletionMessage
)

func InitMCPClient(ctx context.Context, sseAddress string) *client.Client {
	if mcpClient != nil {
		return mcpClient
	}
	fmt.Println("Initializing HTTP client...")
	// Create SSE transport
	sseTransport, err := transport.NewSSE(sseAddress)

	// Create client with the transport
	c := client.NewClient(sseTransport)
	// Start the client
	err = c.Start(ctx)
	CheckErr(err)

	// Initialize the client
	fmt.Println("Initializing client...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "MCP-Go Simple Client Example",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	serverInfo, err := c.Initialize(ctx, initRequest)
	CheckErr(err)
	mcpServerInfo = serverInfo.ServerInfo

	mcpClient = c
	return mcpClient
}

func InitLLMClient() *openai.Client {
	if openaiClient != nil {
		return openaiClient
	}
	openaiClient = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	return openaiClient
}

func BuildSystemPrompt(ctx context.Context) string {
	tools, err := mcpClient.ListTools(ctx, mcp.ListToolsRequest{})
	CheckErr(err)

	toolSchema := make([]string, 0, 3)
	for _, tool := range tools.Tools {
		inputSchema, _ := tool.InputSchema.MarshalJSON()
		toolSchema = append(toolSchema, fmt.Sprintf("- %s\n  %s\n  %s", tool.Name, tool.Description, inputSchema))
	}
	available_tools := []string{
		fmt.Sprintf("## %s", mcpServerInfo.Name),
		"### Available Tools",
	}
	available_tools = append(available_tools, toolSchema...)
	systemPrompt := strings.Replace(systemPromptTemplate, "<$MCP_INFO$>", strings.Join(available_tools, "\n"), -1)
	return systemPrompt
}
