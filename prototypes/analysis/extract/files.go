package extract

import (
	"go/ast"
	"go/token"
	"go/types"
)

type FileVisitor struct {
	Package string
	Imports []Import
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

				for _, s := range gd.Specs {
					is, ok := s.(*ast.ImportSpec)
					if !ok {
						continue
					}
					i := Import{
						Path: is.Path.Value,
						Doc: Doc{
							Comment: is.Doc.Text() + is.Comment.Text(),
							OfQName: is.Path.Value,
						},
					}

					if is.Name != nil {
						i.Name = is.Name.Name
					}

					v.Imports = append(v.Imports, i)
				}

			}
			return v
		}
	default:
		{
			return v
		}
	}
}
