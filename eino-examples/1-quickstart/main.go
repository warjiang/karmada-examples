package main

import "fmt"

func main() {
	model, _ := openai.NewChatModel(ctx, config) // create an invokable LLM instance
	message, _ := model.Generate(ctx, []*Message{
		SystemMessage("you are a helpful assistant."),
		UserMessage("what does the future AI App look like?")})

}
