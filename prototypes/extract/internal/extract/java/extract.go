package java

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/humanbeeng/lepo/prototypes/extract/internal/execute"
	"github.com/humanbeeng/lepo/prototypes/extract/internal/extract"
)

type JavaExtractor struct {
	targetTypes           []extract.ChunkType
	targetTypesToRulesDir map[extract.ChunkType]string
}

func NewJavaExtractor() *JavaExtractor {
	var targetTypes []extract.ChunkType
	targetTypes = append(targetTypes,
		extract.Method,
		extract.Class,
		// TODO: Add interface rule
		extract.Import,
	)

	targetTypesToRulesDir := make(map[extract.ChunkType]string)
	targetTypesToRulesDir[extract.Method] = "internal/extract/java/rules/method.yml"
	targetTypesToRulesDir[extract.Import] = "internal/extract/java/rules/rules/import.yml"
	targetTypesToRulesDir[extract.Constructor] = "internal/extract/java/rules/constructor.yml"
	targetTypesToRulesDir[extract.Import] = "internal/extract/java/rules/import.yml"
	targetTypesToRulesDir[extract.Package] = "internal/extract/java/rules/package.yml"
	targetTypesToRulesDir[extract.Field] = "internal/extract/java/rules/field.yml"

	return &JavaExtractor{
		targetTypes:           targetTypes,
		targetTypesToRulesDir: targetTypesToRulesDir,
	}
}

func (je *JavaExtractor) Extract(file string) ([]extract.Chunk, error) {
	fileinfo, err := os.Stat(file)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}

	if fileinfo.IsDir() {
		err := fmt.Errorf("Error: %v is a directory", file)
		return nil, err
	}

	ext := filepath.Ext(file)
	if ext != ".java" {
		err := fmt.Errorf("Error: %v is not a java file", file)
		return nil, err
	}

	var packageStmt string
	var chunks []extract.Chunk

	for chunkType, rulepath := range je.targetTypesToRulesDir {
		if _, err := os.Stat(rulepath); os.IsNotExist(err) {
			log.Printf("Error: %v rulepath does not exist, hence skipping it", rulepath)
			continue
		}

		cmd := fmt.Sprintf("ast-grep scan -r %v %v --json", rulepath, file)

		stdout, stderr, err := execute.CommandExecute(cmd)

		if stderr != "" {
			log.Printf("Error: %v\n", stderr)
			continue
		}

		if err != nil {
			log.Printf("Error: %v\n", err)
			continue
		}

		var grepResults []extract.GrepResult

		err = json.Unmarshal([]byte(stdout), &grepResults)
		if err != nil {
			log.Printf("Error: Unable to unmarshal stdout %v", err)
			continue
		}

		if chunkType == extract.Package {
			result := grepResults[0]
			packageStmt = result.Text
			log.Printf("Info: %v belongs to package %v", file, packageStmt)
		}

		hasher := sha256.New()

		for _, result := range grepResults {
			hasher.Write([]byte(result.Text))
			contentHash := hex.EncodeToString(hasher.Sum(nil))
			chunk := extract.Chunk{
				File:        result.File,
				Language:    extract.Java,
				Type:        chunkType,
				Content:     result.Text,
				ContentHash: contentHash,
			}
			chunks = append(chunks, chunk)
			hasher.Reset()
		}
	}

	// Populate package hash
	for idx, chunk := range chunks {
		chunk.BelongsTo = packageStmt
		chunks[idx] = chunk
	}

	return chunks, nil
}
