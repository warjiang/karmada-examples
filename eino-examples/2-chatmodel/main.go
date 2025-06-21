package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {

	}
	ctx := context.TODO()
	model := initChatModel(ctx, "openai")

	/*
		streamResult, _ := model.Stream(ctx, []*schema.Message{
			schema.SystemMessage("you are a helpful assistant."),
			schema.UserMessage("what does the future AI App look like?"),
		})
		reportStream(streamResult)
	*/
	generateResult, _ := model.Generate(ctx, []*schema.Message{
		schema.SystemMessage("you are a helpful assistant."),
		schema.UserMessage("what does the future AI App look like?"),
	})
	fmt.Printf("%s", generateResult.Content)
}

func initChatModel(ctx context.Context, modelType string) model.ChatModel {
	var chatModel model.ChatModel
	switch modelType {
	case "ollama":
		config := &ollama.ChatModelConfig{
			BaseURL: os.Getenv("OLLAMA_BASE_URL"),
			Model:   os.Getenv("OLLAMA_MODEL"),
		}
		chatModel, _ = ollama.NewChatModel(ctx, config)
	case "openai":
		config := &openai.ChatModelConfig{
			BaseURL: os.Getenv("OPENAI_BASE_URL"),
			APIKey:  os.Getenv("OPENAI_API_KEY"),
			Model:   os.Getenv("OPENAI_MODEL_NAME"),
		}
		chatModel, _ = openai.NewChatModel(ctx, config)
	case "ark":
		fallthrough
	default:
		config := &ark.ChatModelConfig{
			BaseURL: os.Getenv("ARK_BASE_URL"),
			APIKey:  os.Getenv("ARK_API_KEY"),
			Model:   os.Getenv("ARK_MODEL_NAME"),
		}
		chatModel, _ = ark.NewChatModel(ctx, config)
	}
	return chatModel
}

func reportStream(sr *schema.StreamReader[*schema.Message]) {
	defer sr.Close()

	i := 0
	for {
		message, err := sr.Recv()
		if err == io.EOF { // 流式输出结束
			return
		}
		if err != nil {
			log.Fatalf("recv failed: %v", err)
		}
		// log.Printf("message[%d]: %+v\n", i, message)
		fmt.Printf(message.Content)
		i++
	}
}
