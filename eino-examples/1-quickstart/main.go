package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {

	}
	ctx := context.TODO()
	config := &openai.ChatModelConfig{
		BaseURL: os.Getenv("BASEURL"),
		Model:   os.Getenv("MODEL"),
		APIKey:  os.Getenv("APIKEY"),
	}
	model, _ := openai.NewChatModel(ctx, config) // create an invokable LLM instance
	message, _ := model.Generate(ctx, []*schema.Message{
		schema.SystemMessage("you are a helpful assistant."),
		schema.UserMessage("what does the future AI App look like?"),
	})
	fmt.Println(message.Content)
}
