package embed

import (
	"context"
	"log"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate/entities/models"
)

var WeaviateClient *weaviate.Client

func init() {
	cfg := weaviate.Config{
		Host:       "lepo-test-cluster-wkbl4unf.weaviate.network", // Replace with your endpoint
		Scheme:     "https",
		AuthConfig: auth.ApiKey{Value: "seipFnXH7CxITAJfCjMAr9qBDhevkMCDlztf"}, // Replace w/ your Weaviate instance API key
		Headers: map[string]string{
			"X-OpenAI-Api-Key": "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",
		},
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	className := "CodeSnippets"
	err = client.Schema().AllDeleter().Do(context.TODO())
	if err != nil {
		log.Printf("err: Unable to delete all schema %v\n", err)
		return
	}

	exists, err := client.Schema().ClassExistenceChecker().WithClassName(className).Do(context.Background())

	if exists {
		log.Printf("%v already exists, deleting !\n", className)
		err := client.Schema().ClassDeleter().WithClassName(className).Do(context.Background())
		if err != nil {
			log.Printf("Unable to delete %v\n", className)
			return
		}
	} else {
		log.Println(className, " does not exists")
	}

	WeaviateClient = client
	classObj := &models.Class{
		Class:      className,
		Vectorizer: "text2vec-openai",
		ModuleConfig: map[string]any{
			"text2vec-openai": map[string]any{
				"model":              "ada",
				"modelVersion":       "002",
				"type":               "text",
				"vectorizeClassName": true,
			},
		},
	}

	if client.Schema().ClassCreator().WithClass(classObj).Do(context.Background()) != nil {
		log.Printf("error: %v", err)
		panic(err)
	}
}
