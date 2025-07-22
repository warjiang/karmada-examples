package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create client based on transport type
	var c *client.Client
	var err error

	fmt.Println("Initializing HTTP client...")
	// Create HTTP transport
	sseTransport, err := transport.NewSSE("http://localhost:1234/mcp/sse")
	// NOTE: the default streamableHTTP transport is not 100% identical to the stdio client.
	// By default, it could not receive global notifications (e.g. toolListChanged).
	// You need to enable the `WithContinuousListening()` option to establish a long-live connection,
	// and receive the notifications any time the server sends them.
	//
	//   sseTransport, err := transport.NewStreamableHTTP(*httpURL, transport.WithContinuousListening())
	if err != nil {
		log.Fatalf("Failed to create HTTP transport: %v", err)
	}

	// Create client with the transport
	c = client.NewClient(sseTransport)
	// Start the client
	if err := c.Start(ctx); err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}

	c.OnNotification(func(notification mcp.JSONRPCNotification) {
		fmt.Printf("Received notification: %s\n", notification.Method)
	})

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
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	tools, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	toolSchema := make([]string, 0, 3)
	for _, tool := range tools.Tools {
		inputSchema, _ := tool.InputSchema.MarshalJSON()
		toolSchema = append(toolSchema, fmt.Sprintf("- %s\n  %s\n  %s", tool.Name, tool.Description, inputSchema))
	}
	available_tools := []string{
		fmt.Sprintf("## %s", serverInfo.ServerInfo.Name),
		"### Available Tools",
	}
	available_tools = append(available_tools, toolSchema...)

	//strings.Join(available_tools, "\n")
	content, err := os.ReadFile("/Users/warjiang/workspace/private/karmada-examples/mcp-integration/ex02/prompt.txt")
	if err != nil {
		log.Fatal(err)
	}
	systemPrompt := strings.Replace(string(content), "<$MCP_INFO$>", strings.Join(available_tools, "\n"), -1)
	fmt.Println(systemPrompt)
}
