package extract

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
)

type FileVisitor struct {
	Package string
	Imports map[string][]string
	Fset    *token.FileSet
	Info    *types.Info
}

func (v *FileVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch nd := node.(type) {
	case *ast.ImportSpec:
		{
			fmt.Println("Import:", nd.Path.Value)
		}
	default:
		{
			return v
		}
	}

	return nil
}
