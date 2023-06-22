package sync

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/humanbeeng/lepo/prototypes/extract/internal/embed"
	"github.com/humanbeeng/lepo/prototypes/extract/internal/extract"
	"github.com/humanbeeng/lepo/prototypes/extract/internal/extract/golang"
	"github.com/humanbeeng/lepo/prototypes/extract/internal/extract/java"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

type Syncer interface {
	Sync(path string) error
}

type DirectorySyncer struct {
	languageToExtractor    map[string]extract.Extractor
	ExcludedFolderPatterns []string
}

type DirectorySyncerOpts struct {
	ExcludedFolderPatterns []string
}

func NewDirectorySyncer(opts DirectorySyncerOpts) *DirectorySyncer {
	return &DirectorySyncer{
		languageToExtractor:    buildSupportedLanguagesMap(),
		ExcludedFolderPatterns: opts.ExcludedFolderPatterns,
	}
}

func (s *DirectorySyncer) Sync(path string) error {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return fmt.Errorf("Error: %v does not exists\n", path)
	}

	if !info.IsDir() {
		return fmt.Errorf("Error: %v is not a directory\n", path)
	}

	log.Printf("Sync requested for %v\n", path)

	var dirChunks []extract.Chunk
	err = filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error: %v\n", err)
				return err
			}

			var extractFailedFiles []string

			if !info.IsDir() {
				if extractor, ok := s.languageToExtractor[filepath.Ext(info.Name())]; ok {
					log.Println("Starting extraction for :", path)
					fileChunks, err := extractor.Extract(path)
					if err != nil {
						extractFailedFiles = append(extractFailedFiles, info.Name())
					}
					dirChunks = append(dirChunks, fileChunks...)
				}
			}

			return nil
		})

	fmt.Println("Number of files chunked ", len(dirChunks))

	objects := make([]*models.Object, 0)

	for _, chunk := range dirChunks {
		objects = append(objects, &models.Object{
			Class: "CodeSnippets",
			Properties: map[string]any{
				"file":      chunk.File,
				"code":      chunk.Content,
				"language":  chunk.Language,
				"codeType":  chunk.Type,
				"belongsTo": chunk.BelongsTo,
			},
		})
	}

	log.Printf("Size of objects %v", len(objects))

	batchRes, err := embed.WeaviateClient.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
	if err != nil {
		panic(err)
	}

	for _, res := range batchRes {
		if res.Result.Errors != nil {
			log.Println(res.Result.Errors)
		}
	}
	//
	fields := []graphql.Field{
		{Name: "code"},
		{Name: "language"},
		{Name: "belongsTo"},
		{Name: "file"},
		{Name: "codeType"},
	}

	// openaiClient := openai.NewClient("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")

	// nearVec := embed.WeaviateClient.GraphQL().NearVectorArgBuilder().WithVector(resp.Data[0].Embedding)
	nearText := embed.WeaviateClient.GraphQL().
		NearTextArgBuilder().
		WithConcepts([]string{"health check"})

	result, err := embed.WeaviateClient.GraphQL().Get().
		WithClassName("CodeSnippets").
		WithLimit(2).
		WithNearText(nearText).
		WithFields(fields...).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Data)

	return nil
}

func buildSupportedLanguagesMap() map[string]extract.Extractor {
	supportedLanguages := make(map[string]extract.Extractor)
	supportedLanguages[".go"] = golang.NewGoExtractor()
	supportedLanguages[".java"] = java.NewJavaExtractor()
	return supportedLanguages
}
