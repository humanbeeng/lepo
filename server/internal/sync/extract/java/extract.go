package java

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/humanbeeng/lepo/server/internal/sync/extract"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
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

	lang := golang.GetLanguage()
	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	b, err := ioutil.ReadFile(file) // b has type []byte
	if err != nil {
		err := fmt.Errorf("error: Unable to read %v", file)
		return nil, err
	}

	tree := parser.Parse(nil, b)

	q, err := sitter.NewQuery([]byte(`(
	(package_declaration)? @package
    (import_declaration)? @imports
    
    
    (line_comment)? @line_comment
	(
    	class_declaration
        
        body: (
        	class_body
			(field_declaration)? @field_declaration
            (constructor_declaration)? @constructor
			(	
            	(line_comment)? @method_line_comment
            	(block_comment)? @method_block_comment
            	(method_declaration
            		(throws)? @exception_signature
            	)? @method_declaration
            )  
            (class_declaration)? @inner_class
        )
        
    )? @class_declaration
)`), lang)
	if err != nil {
		return nil, fmt.Errorf("Unable to create query")
	}

	n := tree.RootNode()

	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)

	var funcs []*sitter.Node
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, c := range m.Captures {
			funcs = append(funcs, c.Node)
			fmt.Println("-", funcName(input, c.Node))
		}
	}
	return nil, err
}
