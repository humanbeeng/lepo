package sync

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/humanbeeng/lepo/server/internal/database"
	"github.com/humanbeeng/lepo/server/internal/git"
	"github.com/humanbeeng/lepo/server/internal/sync/extract"
	"github.com/humanbeeng/lepo/server/internal/sync/extract/golang"
	"github.com/humanbeeng/lepo/server/internal/sync/extract/java"
	"github.com/weaviate/weaviate/entities/models"
)

type GitSyncer struct {
	languageToExtractor    map[string]extract.Extractor
	ExcludedFolderPatterns []string
	URL                    string
}

type GitSyncerOpts struct {
	URL string
}

func NewGitSyncer(opts GitSyncerOpts) *GitSyncer {
	return &GitSyncer{
		languageToExtractor:    buildSupportedLanguagesMap(),
		ExcludedFolderPatterns: make([]string, 0),
		URL:                    opts.URL,
	}
}

func (s *GitSyncer) Sync() error {
	// TODO: Introduce goroutines and pass chunks as batch through channel to improve performance

	syncId := uuid.New()

	log.Printf("Sync job requested for %v\n", s.URL)

	// Clone repository
	path := "/home/personal/projects/readonly/cloned/" + syncId.String()
	clonerOpts := git.GitClonerOpts{
		URL:        s.URL,
		TargetPath: path,
	}

	cloner := git.NewGitCloner(clonerOpts)

	err := cloner.Clone()
	if err != nil {
		log.Println("Clone failed", err)
		return err
	}

	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return fmt.Errorf("error: %v does not exists\n", path)
	}

	if !info.IsDir() {
		return fmt.Errorf("error: %v is not a directory\n", path)
	}

	var extractFailedFiles []string
	var dirChunks []extract.Chunk

	err = filepath.Walk(path,

		// TODO: Exclude fileChunks
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("error: %v\n", err)
				return err
			}

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

	if err != nil {
		log.Printf("error: %v\n", err)
	}
	fmt.Println("Number of files chunked", len(dirChunks))
	fmt.Println("Number of files failed", len(extractFailedFiles))

	if len(dirChunks) == 0 {
		fmt.Println("No files found to chunk. Exiting")
		return nil
	}

	// Move this to embed package
	objects := make([]*models.Object, 0)

	for _, chunk := range dirChunks {
		objects = append(objects, &models.Object{
			Class: "CodeSnippets",
			Properties: map[string]any{
				"file":     chunk.File,
				"code":     chunk.Content,
				"language": chunk.Language,
				"codeType": chunk.Type,
				"module":   chunk.Module,
			},
		})
	}

	log.Printf("Number of objects %v\n", len(objects))

	batchRes, err := database.WeaviateClient.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
	if err != nil {
		panic(err)
	}

	for _, res := range batchRes {
		if res.Result.Errors != nil {
			log.Println(res.Result.Errors)
		}
	}
	return nil
}

func buildSupportedLanguagesMap() map[string]extract.Extractor {
	supportedLanguages := make(map[string]extract.Extractor)
	supportedLanguages[".go"] = golang.NewGoExtractor()
	supportedLanguages[".java"] = java.NewJavaExtractor()
	return supportedLanguages
}
