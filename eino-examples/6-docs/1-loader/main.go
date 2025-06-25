package main

import (
	"context"
	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/compose"
	callbacksHelper "github.com/cloudwego/eino/utils/callbacks"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("load dotenv failed")
	}

	accessKey := os.Getenv("OPENAI_API_KEY")

	ctx := context.Background()

	var (
		defaultDim = 1024
	)

	embedder, err := openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
		APIKey:     accessKey,
		Model:      "text-embedding-3-large",
		Dimensions: &defaultDim,
		Timeout:    0,
	})
	if err != nil {
		log.Fatalf("NewEmbedder of openai failed, err=%v", err)
	}

	log.Printf("===== call Embedder directly =====")

	vectors, err := embedder.EmbedStrings(ctx, []string{"hello", "how are you"})
	if err != nil {
		log.Fatalf("EmbedStrings of openai failed, err=%v", err)
	}

	log.Printf("vectors : %v", vectors)

	log.Printf("===== call Embedder in Chain =====")

	handlerHelper := &callbacksHelper.EmbeddingCallbackHandler{
		OnStart: func(ctx context.Context, runInfo *callbacks.RunInfo, input *embedding.CallbackInput) context.Context {
			log.Printf("input access, len: %v, content: %s\n", len(input.Texts), input.Texts)
			return ctx
		},
		OnEnd: func(ctx context.Context, runInfo *callbacks.RunInfo, output *embedding.CallbackOutput) context.Context {
			log.Printf("output finished, len: %v\n", len(output.Embeddings))
			return ctx
		},
	}

	handler := callbacksHelper.NewHandlerHelper().
		Embedding(handlerHelper).
		Handler()

	chain := compose.NewChain[[]string, [][]float64]()
	chain.AppendEmbedding(embedder)

	// 编译并运行
	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("chain Compile failed, err=%v", err)
	}

	vectors, err = runnable.Invoke(ctx, []string{"hello", "how are you"},
		compose.WithCallbacks(handler))
	if err != nil {
		log.Fatalf("Invoke of runnable failed, err=%v", err)
	}

	log.Printf("vectors in chain: %v", vectors)
}
