package main

import (
	"context"
	"eino-examples/common/logs"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
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
	updateTool, err := utils.InferTool("update_todo", "Update a todo item, eg: content,deadline...", UpdateTodoFunc)
	if err != nil {
		logs.Errorf("InferTool failed, err=%v", err)
		return
	}

	// 初始化 tools
	todoTools := []tool.BaseTool{
		getAddTodoTool(), // 使用 NewTool 方式
		updateTool,       // 使用 InferTool 方式
		&ListTodoTool{},  // 使用结构体实现方式, 此处未实现底层逻辑
		//searchTool,
	}

	ctx := context.TODO()
	config := &openai.ChatModelConfig{
		BaseURL: os.Getenv("BASEURL"),
		Model:   os.Getenv("MODEL"),
		APIKey:  os.Getenv("APIKEY"),
	}
	chatModel, _ := openai.NewChatModel(ctx, config) // create an invokable LLM instance

	// 获取工具信息, 用于绑定到 ChatModel
	toolInfos := make([]*schema.ToolInfo, 0, len(todoTools))
	var info *schema.ToolInfo
	for _, todoTool := range todoTools {
		info, err = todoTool.Info(ctx)
		if err != nil {
			logs.Infof("get ToolInfo failed, err=%v", err)
			return
		}
		toolInfos = append(toolInfos, info)
	}

	// 将 tools 绑定到 ChatModel
	err = chatModel.BindTools(toolInfos)
	if err != nil {
		logs.Errorf("BindTools failed, err=%v", err)
		return
	}

	// 创建 tools 节点
	todoToolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: todoTools,
	})
	if err != nil {
		logs.Errorf("NewToolNode failed, err=%v", err)
		return
	}

	// 构建完整的处理链
	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	chain.
		AppendChatModel(chatModel, compose.WithNodeName("chat_model")).
		AppendToolsNode(todoToolsNode, compose.WithNodeName("tools"))

	// 编译并运行 chain
	agent, err := chain.Compile(ctx)
	if err != nil {
		logs.Errorf("chain.Compile failed, err=%v", err)
		return
	}

	// 运行示例
	resp, err := agent.Invoke(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "添加一个学习 Eino 的 TODO，同时搜索一下 cloudwego/eino 的仓库地址",
		},
	})
	if err != nil {
		logs.Errorf("agent.Invoke failed, err=%v", err)
		return
	}

	// 输出结果
	for idx, msg := range resp {
		logs.Infof("\n")
		logs.Infof("message %d: %s: %s", idx, msg.Role, msg.Content)
	}

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
