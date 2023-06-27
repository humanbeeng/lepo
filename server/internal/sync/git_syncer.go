package sync

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/humanbeeng/lepo/server/internal/git"
	"github.com/humanbeeng/lepo/server/internal/sync/extract"
	"github.com/humanbeeng/lepo/server/internal/sync/extract/golang"
	"github.com/humanbeeng/lepo/server/internal/sync/extract/java"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
	"go.uber.org/zap"
)

type GitSyncer struct {
	languageToExtractor    map[string]extract.Extractor
	ExcludedFolderPatterns []string
	cloner                 *git.GitCloner
	weaviate               *weaviate.Client
	logger                 *zap.Logger
}

type GitSyncerOpts struct {
	URL string
}

func NewGitSyncer(cloner *git.GitCloner, weaviate *weaviate.Client, logger *zap.Logger) *GitSyncer {
	return &GitSyncer{
		languageToExtractor:    buildSupportedLanguagesMap(logger),
		ExcludedFolderPatterns: make([]string, 0),
		cloner:                 cloner,
		weaviate:               weaviate,
		logger:                 logger,
	}
}

func (s *GitSyncer) Sync(url string) error {
	// TODO: Introduce goroutines and pass chunks as batch through channel to improve performance

	syncId := uuid.New()
	targetPath := filepath.Join("../temp/repo", syncId.String())

	s.logger.Info("Sync job requested", zap.String("url", url))

	// Clone repository
	req := git.GitCloneRequest{
		URL:        url,
		TargetPath: targetPath,
	}

	err := s.cloner.Clone(req)
	if err != nil {
		s.logger.Error("Clone failed", zap.Error(err))
		return err
	}

	info, err := os.Stat(targetPath)

	if os.IsNotExist(err) {
		return fmt.Errorf("error: %v does not exists\n", url)
	}

	if !info.IsDir() {
		return fmt.Errorf("error: %v is not a directory\n", url)
	}

	defer s.cleanup(targetPath)

	var extractFailedFiles []string
	var dirChunks []extract.Chunk

	err = filepath.Walk(targetPath,

		// TODO: Exclude fileChunks
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error: %v\n", err)
			}

			if !info.IsDir() {
				if extractor, ok := s.languageToExtractor[filepath.Ext(info.Name())]; ok {
					s.logger.Info("Starting extraction", zap.String("path", path))
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
		s.logger.Error("Error while walking", zap.Error(err))
	}

	s.logger.Info("Number of files chunked", zap.Int("chunked", len(dirChunks)))
	s.logger.Info("Number of files failed", zap.Int("failed", len(extractFailedFiles)))

	if len(dirChunks) == 0 {
		s.logger.Warn("No files found to chunk. Exiting")
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

	s.logger.Info("Pushing to weaviate", zap.Int("objects", len(objects)))

	batchRes, err := s.weaviate.Batch().
		ObjectsBatcher().
		WithObjects(objects...).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	for _, res := range batchRes {
		if res.Result.Errors != nil {
			s.logger.Warn("Error while pushing", zap.Any("errors", res.Result.Errors))
		}
	}

	s.logger.Info("Sync completed", zap.String("syncId", syncId.String()))
	return nil
}

func (s *GitSyncer) Desync() error {
	return s.weaviate.Schema().AllDeleter().Do(context.Background())
}

func buildSupportedLanguagesMap(logger *zap.Logger) map[string]extract.Extractor {
	supportedLanguages := make(map[string]extract.Extractor)
	supportedLanguages[".go"] = golang.NewGoExtractor(logger)
	supportedLanguages[".java"] = java.NewJavaExtractor(logger)
	return supportedLanguages
}

func (s *GitSyncer) cleanup(targetPath string) {
	err := os.RemoveAll(targetPath)
	s.logger.Info("Cleaning up", zap.String("targetPath", targetPath))
	if err != nil {
		log.Println("warn: Unable to cleanup cloned repo", targetPath)
	}
}
