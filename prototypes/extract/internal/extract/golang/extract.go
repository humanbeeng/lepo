package golang

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

type GoExtractor struct {
	targetTypes           []extract.ChunkType
	targetTypesToRulesDir map[extract.ChunkType]string
}

func NewGoExtractor() *GoExtractor {
	var targetTypes []extract.ChunkType

	targetTypes = append(
		targetTypes,
		extract.Method,
		extract.Interface,
		extract.Struct,
		extract.Import,
		extract.Function,
	)

	targetTypesToRulesDir := make(map[extract.ChunkType]string)

	targetTypesToRulesDir[extract.Struct] = "internal/extract/golang/rules/struct.yml"
	targetTypesToRulesDir[extract.Method] = "internal/extract/golang/rules/method.yml"
	targetTypesToRulesDir[extract.Function] = "internal/extract/golang/rules/function.yml"
	targetTypesToRulesDir[extract.Import] = "internal/extract/golang/rules/import.yml"
	targetTypesToRulesDir[extract.Package] = "internal/extract/golang/rules/package.yml"

	return &GoExtractor{
		targetTypes:           targetTypes,
		targetTypesToRulesDir: targetTypesToRulesDir,
	}
}

func (ge *GoExtractor) Extract(file string) ([]extract.Chunk, error) {
	// TODO: Move this to util package
	fileinfo, err := os.Stat(file)
	if err != nil {
		log.Printf("Error: %v does not exist %v", file, err)
		return nil, err
	}

	if fileinfo.IsDir() {
		err := fmt.Errorf("Error: %v is a directory", file)
		return nil, err
	}

	ext := filepath.Ext(file)
	if ext != ".go" {
		err := fmt.Errorf("Error: %v is not a Go file", file)
		return nil, err
	}

	var packageStmt string
	var chunks []extract.Chunk

	for chunkType, rulepath := range ge.targetTypesToRulesDir {
		if _, err := os.Stat(rulepath); os.IsNotExist(err) {
			log.Printf("Error: %v rulepath does not exist", rulepath)
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
			log.Printf("Error: Unable to unmarshal json %v\n", err)
			continue
		}

		if len(grepResults) == 0 {
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
				File:        file,
				Language:    extract.Go,
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
