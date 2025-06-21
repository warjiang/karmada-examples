package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func useTemplateWithGraph() {
	// 创建模板
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个{role}。"),
		schema.MessagesPlaceholder("history_key", false),
		&schema.Message{
			Role:    schema.User,
			Content: "请帮我{task}。",
		},
	)

	graph := compose.NewGraph[map[string]any, []*schema.Message]()
	graph.AddChatTemplateNode("template_node", template)

	ctx := context.TODO()

	variables := map[string]any{
		"role": "专业的助手",
		"task": "写一首诗",
		"history_key": []*schema.Message{
			{Role: schema.User, Content: "告诉我油画是什么?"},
			{Role: schema.Assistant, Content: "油画是xxx"},
		},
	}

	graph.AddLambdaNode("your_node_key", compose.InvokableLambda(func(ctx context.Context, input map[string]any) (x []*schema.Message, err error) {
		// your logic
		return
	}))
	graph.AddEdge(compose.START, "your_node_key")
	graph.AddEdge("your_node_key", compose.END)

	runnable, err := graph.Compile(ctx)
	if err != nil {
		panic(err)
	}

	result, err := runnable.Invoke(ctx, variables)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", result)
}
