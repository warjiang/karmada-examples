package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/beevik/etree"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
	"regexp"
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
	content, err := os.ReadFile("/Users/warjiang/workspace/private/karmada-examples/mcp-integration/ex03/prompt.txt")
	if err != nil {
		log.Fatal(err)
	}
	systemPrompt := strings.Replace(string(content), "<$MCP_INFO$>", strings.Join(available_tools, "\n"), -1)

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	var msgs = []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "How many member clusters are there in the karmada control-plane?",
		},
	}

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
		Messages: msgs,
	}
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	//fmt.Println(resp.Choices[0].Message.Content)
	_, toolName, toolArgs, err := ParseToolString(resp.Choices[0].Message.Content)
	if err != nil {
		fmt.Printf("ParseToolString error: %v\n", err)
		return
	}

	//fmt.Printf("ServerName: %s\n", serverName)
	//fmt.Printf("toolName: %s\n", toolName)
	//fmt.Printf("toolArgs: %s\n", toolArgs)

	toolResp, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      toolName,
			Arguments: toolArgs,
			Meta:      nil,
		},
	})
	if err != nil {
		fmt.Printf("CallTool error: %v\n", err)
		return
	}
	marshalJSON, _ := toolResp.MarshalJSON()
	//fmt.Println(string(marshalJSON))

	textContent, ok := toolResp.Content[0].(mcp.TextContent)
	if !ok {
		fmt.Println("Content is not of type TextContent")
		return
	}

	msgs = append(msgs, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: textContent.Text,
	})

	msgs = append(msgs,
		//"role": "user",
		//	"content": f"[Tool {tool_name} \n returned: {result}]"
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("[Tool %s \n returned: %s]", toolName, marshalJSON),
		},
	)
	//fmt.Println("finished")

	req = openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
		Messages: msgs,
	}
	resp, err = client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	fmt.Printf("Bot: %s\n", resp.Choices[0].Message.Content)
}

// ParseToolString 解析大模型工具调用返回的字符串
func ParseToolString(toolString string) (string, string, map[string]interface{}, error) {
	// 使用正则表达式提取 XML 部分
	re := regexp.MustCompile(`(?s)<use_mcp_tool>.*?</use_mcp_tool>`)
	matches := re.FindStringSubmatch(toolString)
	if len(matches) == 0 {
		return "", "", nil, errors.New("no tool XML found in input string")
	}
	xmlStr := matches[0]

	// 解析 XML
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xmlStr); err != nil {
		return "", "", nil, fmt.Errorf("failed to parse XML: %v", err)
	}

	root := doc.SelectElement("use_mcp_tool")
	if root == nil {
		return "", "", nil, errors.New("invalid XML structure: missing use_mcp_tool element")
	}

	// 提取 server_name
	serverNameElem := root.SelectElement("server_name")
	if serverNameElem == nil {
		return "", "", nil, errors.New("missing server_name element")
	}
	serverName := serverNameElem.Text()

	// 提取 tool_name
	toolNameElem := root.SelectElement("tool_name")
	if toolNameElem == nil {
		return "", "", nil, errors.New("missing tool_name element")
	}
	toolName := toolNameElem.Text()

	// 提取并解析 arguments
	argsElem := root.SelectElement("arguments")
	if argsElem == nil {
		return "", "", nil, errors.New("missing arguments element")
	}

	var toolArgs map[string]interface{}
	argsText := strings.TrimSpace(argsElem.Text())
	if argsText != "" {
		if err := json.Unmarshal([]byte(argsText), &toolArgs); err != nil {
			return "", "", nil, fmt.Errorf("invalid tool arguments: %v", err)
		}
	} else {
		toolArgs = make(map[string]interface{})
	}

	return serverName, toolName, toolArgs, nil
}
