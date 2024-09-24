package golang

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/humanbeeng/lepo/prototypes/analysis/extract"
)

type VarVisitor struct {
	Vars map[string]extract.Variable
	Fset *token.FileSet
	Info *types.Info
}

func (v *VarVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch nd := node.(type) {
	case *ast.GenDecl:
		{
			for _, spec := range nd.Specs {
				vs, vsOk := spec.(*ast.ValueSpec)
				if !vsOk {
					return v
				}
				for _, vsName := range vs.Names {
					vsObj, objOk := v.Info.Defs[vsName]
					if !objOk {
						return v
					}

					name := vsObj.Name()
					namespace := extract.Namespace{Name: vsObj.Pkg().Path()}
					qname := namespace.Name + "." + name
					pos := v.Fset.Position(vs.Pos()).Line
					end := v.Fset.Position(vs.End()).Line
					filepath := v.Fset.Position(vs.Pos()).Filename

					variable := extract.Variable{
						Name:       name,
						QName:      qname,
						TypeQName:  vsObj.Type().String(),
						Underlying: vsObj.Type().Underlying().String(),
						Namespace:  namespace,
						Pos:        pos,
						End:        end,
						Filepath:   filepath,
						Doc: extract.Doc{
							Comment: vs.Doc.Text() + vs.Comment.Text(),
							OfQName: qname,
							// TODO: Revisit
							Type: extract.SingleLine,
						},
					}

					v.Vars[qname] = variable

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
