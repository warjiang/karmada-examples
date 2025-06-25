package main

import (
	"context"
	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("load dotenv failed")
	}
	// Get the environment variables
	addr := os.Getenv("MILVUS_ADDR")
	username := os.Getenv("MILVUS_USERNAME")
	password := os.Getenv("MILVUS_PASSWORD")
	//arkApiKey := os.Getenv("ARK_API_KEY")
	//arkModel := os.Getenv("ARK_MODEL")

	// Create a client
	ctx := context.Background()
	cli, err := client.NewClient(ctx, client.Config{
		Address:  addr,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return
	}
	defer cli.Close()

	/*
		// Create an embedding model
		emb, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
			APIKey: arkApiKey,
			Model:  arkModel,
		})
	*/
	//var (
	//	defaultDim = 1024
	//)
	emb, err := openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
		APIKey: os.Getenv("OPENAI_API_KEY"),
		Model:  "text-embedding-3-large",
		//Dimensions: &defaultDim,
		Timeout: 0,
	})
	if err != nil {
		log.Fatalf("Failed to create embedding: %v", err)
		return
	}

	// Create an indexer
	idxer, err := milvus.NewIndexer(ctx, &milvus.IndexerConfig{
		Client:     cli,
		Embedding:  emb,
		Collection: "customized_setup_2",
	})
	if err != nil {
		log.Fatalf("Failed to create indexer: %v", err)
		return
	}
	log.Printf("Indexer created success")
	/*
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
	*/
	// Store documents
	docs := []*schema.Document{
		{
			ID:      "milvus-1",
			Content: "milvus is an open-source vector database",
			//MetaData: map[string]any{
			//	"h1": "milvus",
			//	"h2": "open-source",
			//	"h3": "vector database",
			//},
		},
		{
			ID:      "milvus-2",
			Content: "milvus is a distributed vector database",
		},
	}
	ids, err := idxer.Store(ctx, docs)
	if err != nil {
		log.Fatalf("Failed to store: %v", err)
		return
	}
	log.Printf("Store success, ids: %v", ids)
}
