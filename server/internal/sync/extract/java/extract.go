package java

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/humanbeeng/lepo/server/internal/sync/execute"
	"github.com/humanbeeng/lepo/server/internal/sync/extract"
	"go.uber.org/zap"
)

type JavaExtractor struct {
	targetTypes           []extract.ChunkType
	targetTypesToRulesDir map[extract.ChunkType]string
	logger                *zap.Logger
}

func NewJavaExtractor(logger *zap.Logger) *JavaExtractor {
	var targetTypes []extract.ChunkType
	targetTypes = append(targetTypes,
		extract.Method,
		extract.Class,
		// TODO: Add interface rule
		extract.Import,
	)

	targetTypesToRulesDir := make(map[extract.ChunkType]string)
	targetTypesToRulesDir[extract.Method] = "internal/sync/extract/java/rules/method.yml"
	targetTypesToRulesDir[extract.Import] = "internal/sync/extract/java/rules/rules/import.yml"
	targetTypesToRulesDir[extract.Constructor] = "internal/sync/extract/java/rules/constructor.yml"
	targetTypesToRulesDir[extract.Import] = "internal/sync/extract/java/rules/import.yml"
	targetTypesToRulesDir[extract.Package] = "internal/sync/extract/java/rules/package.yml"
	targetTypesToRulesDir[extract.Field] = "internal/sync/extract/java/rules/field.yml"

	return &JavaExtractor{
		targetTypes:           targetTypes,
		targetTypesToRulesDir: targetTypesToRulesDir,
		logger:                logger,
	}
}

func (je *JavaExtractor) Extract(file string) ([]extract.Chunk, error) {
	fileinfo, err := os.Stat(file)
	if err != nil {
		je.logger.Error("error:", zap.Error(err))
		return nil, err
	}

	if fileinfo.IsDir() {
		err := fmt.Errorf("error: %v is a directory", file)
		return nil, err
	}

	ext := filepath.Ext(file)
	if ext != ".java" {
		err := fmt.Errorf("error: %v is not a java file", file)
		return nil, err
	}

	var packageStmt string
	var chunks []extract.Chunk

	for chunkType, rulepath := range je.targetTypesToRulesDir {
		if _, err := os.Stat(rulepath); os.IsNotExist(err) {
			err = fmt.Errorf("%v rulepath does not exist", rulepath)
			return nil, err
		}

		cmd := fmt.Sprintf("ast-grep scan -r %v %v --json", rulepath, file)

		stdout, stderr, err := execute.CommandExecute(cmd)

		if stderr != "" {
			je.logger.Error("error:", zap.String("stderr", stderr))
			continue
		}

		if err != nil {
			je.logger.Error("error:", zap.Error(err))
			continue
		}

		var grepResults []extract.GrepResult

		err = json.Unmarshal([]byte(stdout), &grepResults)
		if err != nil {
			je.logger.Error("Unable to unmarshal stdout", zap.Error(err))
			continue
		}

		if len(grepResults) == 0 {
			continue
		}

		if chunkType == extract.Package {
			result := grepResults[0]
			packageStmt = result.Text
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

	// Populate package name to all chunks
	for idx, chunk := range chunks {
		chunk.Module = packageStmt
		chunks[idx] = chunk
	}

	return chunks, nil
}
