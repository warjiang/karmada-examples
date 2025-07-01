package main

import (
	"bufio"
	"context"
	"github.com/cloudwego/eino/flow/agent/multiagent/host"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"io"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("load dotenv failed")
	}

	/*
		openAIAPIKey := os.Getenv("OPENAI_API_KEY")
		openAIBaseURL := os.Getenv("OPENAI_BASE_URL")
		openAIModelName := os.Getenv("OPENAI_MODEL_NAME")
	*/
	openAIAPIKey := os.Getenv("APIKEY")
	openAIBaseURL := os.Getenv("BASEURL")
	openAIModelName := os.Getenv("MODEL")

	ollamaBaseURL := os.Getenv("OLLAMA_BASE_URL")
	ollamaModel := os.Getenv("OLLAMA_MODEL")

	ctx := context.Background()

	h, err := newHost(ctx, openAIBaseURL, openAIAPIKey, openAIModelName)
	if err != nil {
		panic(err)
	}

	writer, err := newWriteJournalSpecialist(ctx, ollamaBaseURL, ollamaModel)
	if err != nil {
		panic(err)
	}

	reader, err := newReadJournalSpecialist(ctx)
	if err != nil {
		panic(err)
	}

	answerer, err := newAnswerWithJournalSpecialist(ctx, ollamaBaseURL, ollamaModel)
	if err != nil {
		panic(err)
	}

	hostMA, err := host.NewMultiAgent(ctx, &host.MultiAgentConfig{
		Host: *h,
		Specialists: []*host.Specialist{
			writer,
			reader,
			answerer,
		},
	})

	if err != nil {
		panic(err)
	}

	cb := &logCallback{}
	for { // 多轮对话，除非用户输入了 "exit"，否则一直循环
		println("\n\nYou: ") // 提示轮到用户输入了

		var message string
		scanner := bufio.NewScanner(os.Stdin) // 获取用户在命令行的输入
		for scanner.Scan() {
			message += scanner.Text()
			break
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		if message == "exit" {
			return
		}

		msg := &schema.Message{
			Role:    schema.User,
			Content: message,
		}

		out, err := hostMA.Stream(ctx, []*schema.Message{msg}, host.WithAgentCallbacks(cb))
		if err != nil {
			panic(err)
		}

		println("\nAnswer:")

		for {
			msg, err := out.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
			}

			print(msg.Content)
		}

		out.Close()
	}
}

type logCallback struct{}

func (l *logCallback) OnHandOff(ctx context.Context, info *host.HandOffInfo) context.Context {
	println("\nHandOff to", info.ToAgentName, "with argument", info.Argument)
	return ctx
}
