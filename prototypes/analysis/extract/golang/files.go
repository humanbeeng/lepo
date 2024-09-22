package golang

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/humanbeeng/lepo/prototypes/analysis/extract"
)

type FileVisitor struct {
	Package string
	Files   map[string]extract.File
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
			// TODO: Remove absolute paths
			filename := v.Fset.Position(nd.Pos()).Filename

			f := extract.File{
				Filename:  filename,
				Namespace: v.Package,
				Language:  extract.Go,
				Imports:   make([]extract.Import, 0),
			}

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
					i := extract.Import{
						Path: strings.Trim(is.Path.Value, "\""),
						Doc: extract.Doc{
							Comment: is.Doc.Text() + is.Comment.Text(),
							OfQName: is.Path.Value,
						},
					}

					if is.Name != nil {
						i.Name = is.Name.Name
					}
					f.Imports = append(f.Imports, i)
				}
				v.Files[f.Namespace+"."+f.Filename] = f
			}
			return v
		}
	default:
		{
			return v
		}
	}
}
