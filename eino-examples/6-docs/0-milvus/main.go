package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/milvusclient"

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

	// Create a client
	ctx := context.Background()
	cli, err := milvusclient.New(ctx, &milvusclient.ClientConfig{
		Address:  addr,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return
	}
	defer cli.Close(ctx)

	collections, err := cli.ListCollections(ctx, milvusclient.NewListCollectionOption())
	if err != nil {
		panic(err)
	}

	for _, collection := range collections {
		log.Println(collection)
	}
	/*
		err = cli.DropCollection(ctx, milvusclient.NewDropCollectionOption("eino_collection"))
		if err != nil {
			panic(err)
		}
	*/
	/*
		indexOptions := []milvusclient.CreateIndexOption{
			milvusclient.NewCreateIndexOption(collectionName, "my_vector", index.NewAutoIndex(entity.COSINE)),
			milvusclient.NewCreateIndexOption(collectionName, "my_id", index.NewAutoIndex(entity.COSINE)),
		}
	*/

	/*
		collectionName := "customized_setup_1"
		cli.CreateCollection(ctx,
			milvusclient.
				SimpleCreateCollectionOptions(collectionName, 4096).
				WithVarcharPK(true, 64).
				WithShardNum(1),
		)
	*/
	collectionName := "customized_setup_2"
	schema := entity.NewSchema().WithDynamicFieldEnabled(true)
	fields := []*entity.Field{
		entity.NewField().
			WithName(defaultCollectionID).
			WithDescription(defaultCollectionIDDesc).
			WithIsPrimaryKey(true).
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(255),
		entity.NewField().
			WithName(defaultCollectionVector).
			WithDescription(defaultCollectionVectorDesc).
			WithIsPrimaryKey(false).
			WithDataType(entity.FieldTypeBinaryVector).
			WithDim(defaultDim),
		entity.NewField().
			WithName(defaultCollectionContent).
			WithDescription(defaultCollectionContentDesc).
			WithIsPrimaryKey(false).
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(1024),
		entity.NewField().
			WithName(defaultCollectionMetadata).
			WithDescription(defaultCollectionMetadataDesc).
			WithIsPrimaryKey(false).
			WithDataType(entity.FieldTypeJSON),
	}
	for _, field := range fields {
		schema = schema.WithField(field)
	}

	err = cli.CreateCollection(ctx,
		milvusclient.
			NewCreateCollectionOption(
				collectionName,
				schema,
			).
			WithVarcharPK(true, 64).
			WithShardNum(1),
	)
	if err != nil {
		panic(err)
	}
}
