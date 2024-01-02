package extract

import (
	"go/ast"
	"go/token"
	"go/types"
)

type FileVisitor struct {
	Package string
	Imports []string
	Fset    *token.FileSet
	Info    *types.Info
}

func (v *FileVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch nd := node.(type) {
	case *ast.File:
		{

			for _, d := range nd.Decls {
				gd, ok := d.(*ast.GenDecl)
				if !ok {
					continue
				}
				var imports []string

				for _, s := range gd.Specs {
					is, ok := s.(*ast.ImportSpec)
					if !ok {
						continue
					}
					imports = append(imports, is.Path.Value)
				}

				v.Imports = imports
			}
			return v
		}
	default:
		{
			return v
		}
	}
}
