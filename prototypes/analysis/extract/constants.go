package extract

import (
	"go/ast"
	"go/token"
	"go/types"
)

type ConstVisitor struct {
	ast.Visitor
	Constants map[string]Constant
	Fset      *token.FileSet
	Info      *types.Info
}

func (v *ConstVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch nd := node.(type) {
	case *ast.File:
		{
			return v
		}
	case *ast.GenDecl:
		{
			for _, spec := range nd.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				for _, name := range vs.Names {
					vsObj := v.Info.Defs[name]

					pos := v.Fset.Position(vs.Pos()).Line
					end := v.Fset.Position(vs.End()).Line
					filepath := v.Fset.Position(vs.Pos()).Filename

					qname := vsObj.Pkg().Path() + "." + name.Name

					constant := Constant{
						Name:       name.Name,
						QName:      qname,
						TypeQName:  vsObj.Type().String(),
						Underlying: vsObj.Type().Underlying().String(),
						// TODO: Get code
						Code:     "",
						Doc:      Doc{Comment: vs.Doc.Text() + vs.Comment.Text(), OfQName: qname},
						Pos:      pos,
						End:      end,
						Filepath: filepath,
					}
					v.Constants[qname] = constant
				}

			}
		}
	}
	return nil
}
