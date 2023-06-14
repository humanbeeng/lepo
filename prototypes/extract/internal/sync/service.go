package sync

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/humanbeeng/lepo/prototypes/extract/internal/extract"
	"github.com/humanbeeng/lepo/prototypes/extract/internal/extract/golang"
	"github.com/humanbeeng/lepo/prototypes/extract/internal/extract/java"
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

	for _, chunk := range dirChunks {
		fmt.Printf("-------------\n\n%+v\n---------------\n\n\n\n\n", chunk)
	}

	return nil
}

func buildSupportedLanguagesMap() map[string]extract.Extractor {
	supportedLanguages := make(map[string]extract.Extractor)
	supportedLanguages[".go"] = golang.NewGoExtractor()
	supportedLanguages[".java"] = java.NewJavaExtractor()
	return supportedLanguages
}
