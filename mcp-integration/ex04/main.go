package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sashabaranov/go-openai"
	"os"
	"strings"
)

func main() {
	ctx := context.Background()

	reader := bufio.NewReader(os.Stdin)
	InitLLMClient()
	InitMCPClient(ctx, os.Getenv("MCP_SERVER_URL"))
	for {
		fmt.Print("Human: ")
		input, _ := reader.ReadString('\n')
		cmd := strings.TrimSpace(input)

		switch cmd {
		case "exit":
			fmt.Println("Prepare to quit program...")
			return
		case "help":
			fmt.Println("Available commands: help, exit, version")
		case "version":
			fmt.Println("Version 1.0.0")
		default:
			output, err := handle(ctx, cmd)
			if err != nil {
				fmt.Printf("error occurred: %v", err)
				return
			} else {
				fmt.Println(output)
			}
		}
	}
}

func handle(ctx context.Context, input string) (string, error) {
	systemPrompt := BuildSystemPrompt(ctx)
	//fmt.Println(systemPrompt)
	messages = append(messages,
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		},
	)
	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
		Messages: messages,
	}
	resp, err := openaiClient.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	_, toolName, toolArgs, err := ParseToolString(resp.Choices[0].Message.Content)
	if err != nil {
		return resp.Choices[0].Message.Content, nil
	}

	toolResp, err := mcpClient.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      toolName,
			Arguments: toolArgs,
			Meta:      nil,
		},
	})
	if err != nil {
		return "", err
	}

	textContent, ok := toolResp.Content[0].(mcp.TextContent)
	if !ok {
		return "", errors.New("Content is not of type TextContent")
	}

	marshalJSON, _ := toolResp.MarshalJSON()
	messages = append(messages,
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: textContent.Text,
		},
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("[Tool %s \n returned: %s]", toolName, marshalJSON),
		},
	)

	req = openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
		Messages: messages,
	}
	resp, err = openaiClient.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Bot: %s\n", resp.Choices[0].Message.Content), nil
}
