package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func main() {
	// 创建模板
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个{role}。"),
		schema.MessagesPlaceholder("history_key", false),
		&schema.Message{
			Role:    schema.User,
			Content: "请帮我{task}。",
		},
	)

	// 准备变量
	variables := map[string]any{
		"role": "专业的助手",
		"task": "写一首诗",
		"history_key": []*schema.Message{
			{Role: schema.User, Content: "告诉我油画是什么?"},
			{Role: schema.Assistant, Content: "油画是xxx"},
		},
	}

	// 格式化模板
	messages, err := template.Format(context.Background(), variables)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", messages)
}
