package sync

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/humanbeeng/lepo/server/internal/sync/extract"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
	"go.uber.org/zap"
)

// TODO: Delete this
type DirSyncer struct {
	weaviate            *weaviate.Client
	languageToExtractor map[string]extract.Extractor
	logger              *zap.Logger
}

func NewDirSyncer(weaviate *weaviate.Client, logger *zap.Logger) *DirSyncer {
	return &DirSyncer{
		weaviate:            weaviate,
		languageToExtractor: buildSupportedLanguagesMap(logger),
		logger:              logger,
	}
}

func (d *DirSyncer) Sync(targetPath string) error {
	d.logger.Info("Directory sync job requested")

	info, err := os.Stat(targetPath)

	if os.IsNotExist(err) {
		return fmt.Errorf("error: %v does not exists\n", targetPath)
	}

	if !info.IsDir() {
		return fmt.Errorf("error: %v is not a directory\n", targetPath)
	}

	var extractFailedFiles []string
	var dirChunks []extract.Chunk

	err = filepath.Walk(targetPath,

		// TODO: Exclude fileChunks
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error: %v\n", err)
			}

			if !info.IsDir() {
				if extractor, ok := d.languageToExtractor[filepath.Ext(info.Name())]; ok {
					d.logger.Info("Starting extraction", zap.String("path", path))
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
		d.logger.Error("Error while walking", zap.Error(err))
	}

	d.logger.Info("Number of files chunked", zap.Int("chunked", len(dirChunks)))
	d.logger.Info("Number of files failed", zap.Int("failed", len(extractFailedFiles)))

	if len(dirChunks) == 0 {
		d.logger.Warn("No files found to chunk. Exiting")
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

	d.logger.Info("Pushing to weaviate", zap.Int("objects", len(objects)))

	batchRes, err := d.weaviate.Batch().
		ObjectsBatcher().
		WithObjects(objects...).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	for _, res := range batchRes {
		if res.Result.Errors != nil {
			d.logger.Warn("Error while pushing", zap.Any("errors", res.Result.Errors))
		}
	}

	d.logger.Info("Sync completed")
	return nil
}

func (s *DirSyncer) Desync() error {
	resp, err := s.weaviate.Batch().
		ObjectsBatchDeleter().
		WithClassName("CodeSnippets").
		Do(context.Background())
	if err != nil {
		return err
	}

	s.logger.Warn("Failed to delete", zap.Int64("num_failed", resp.Results.Failed))

	return nil
}
