package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/openai"
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
	config := &openai.ChatModelConfig{
		BaseURL: os.Getenv("BASEURL"),
		Model:   os.Getenv("MODEL"),
		APIKey:  os.Getenv("APIKEY"),
	}
	model, _ := openai.NewChatModel(ctx, config) // create an invokable LLM instance
	streamResult, _ := model.Stream(ctx, []*schema.Message{
		schema.SystemMessage("you are a helpful assistant."),
		schema.UserMessage("what does the future AI App look like?"),
	})
	reportStream(streamResult)

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
