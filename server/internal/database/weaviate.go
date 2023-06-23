package database

import (
	"context"
	"fmt"
	"log"

	"github.com/humanbeeng/lepo/server/internal/config"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate/entities/models"
)

func BootStrapWeaviate() (*weaviate.Client, error) {
	log.Println("Initializing Weaviate Client")
	appConfig, err := config.GetAppConfig()
	if err != nil {
		return nil, err
	}
	cfg := weaviate.Config{
		Host:       appConfig.WeaviateConfig.Host,
		Scheme:     appConfig.WeaviateConfig.Scheme,
		AuthConfig: auth.ApiKey{Value: appConfig.WeaviateConfig.ApiKey},
		Headers: map[string]string{
			"X-OpenAI-Api-Key": appConfig.OpenAIConfig.ApiKey,
		},
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	className := "CodeSnippets"

	// // Note: Use status code from ClassDeleter 400 to determine if class exists or not
	// exists, err := client.Schema().ClassExistenceChecker().WithClassName(className).Do(context.Background())

	// // TODO: Revisit and refactor
	// if exists {
	// 	log.Printf("%v already exists, deleting !\n", className)
	// 	err := client.Schema().ClassDeleter().WithClassName(className).Do(context.Background())
	// 	if err != nil {
	// 		log.Printf("Unable to delete %v\n", className)
	// 		return nil, err
	// 	}
	// } else {
	// 	log.Println(className, " does not exists")
	// }

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

	exists, err := client.Schema().
		ClassExistenceChecker().
		WithClassName(className).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exists {
		if err := client.Schema().ClassCreator().WithClass(classObj).Do(context.Background()); err != nil {
			log.Printf("error: Unable to create class %v", err)
			panic(err)
		}
	}

	// Check weaviate status
	readyCheckerRequest := client.Misc().ReadyChecker()
	ready, err := readyCheckerRequest.Do(context.Background())
	if err != nil {
		return nil, err
	}

	if !ready {
		err := fmt.Errorf("error: Weaviate not ready")
		return nil, err
	}

	log.Println("info: Weaviate connection established")

	return client, nil
}
