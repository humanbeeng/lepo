package golang

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/humanbeeng/lepo/prototypes/analysis/extract"
)

type InterfaceVisitor struct {
	ast.Visitor
	Fset       *token.FileSet
	Info       *types.Info
	Interfaces map[string]extract.TypeDecl
	Members    map[string]extract.Member
}

func (v *InterfaceVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch n := node.(type) {

	case *ast.GenDecl:
		{
			for _, s := range n.Specs {
				if ts, ok := s.(*ast.TypeSpec); ok {
					tsObj := v.Info.Defs[ts.Name]
					if inf, ok := ts.Type.(*ast.InterfaceType); ok {
						infQname := tsObj.Pkg().Path() + "." + tsObj.Name()
						pos := v.Fset.Position(inf.Pos()).Line
						end := v.Fset.Position(inf.End()).Line
						filepath := v.Fset.Position(inf.Pos()).Filename

						infCode, err := extractCode(n, v.Fset)
						if err != nil {
							// TODO -- Better error handling
							panic(err)
						}

						td := extract.TypeDecl{
							Name:  ts.Name.Name,
							QName: infQname,
							Namespace: extract.Namespace{
								Name: tsObj.Pkg().Path(),
							},
							TypeQName:  tsObj.Type().String(),
							Underlying: tsObj.Type().Underlying().String(),
							Kind:       extract.Interface,
							Doc: extract.Doc{
								Comment: ts.Doc.Text(),
								OfQName: infQname,
								// TODO: Add comment type
							},
							Pos:      pos,
							End:      end,
							Filepath: filepath,
							Code:     infCode,
						}
						v.Interfaces[infQname] = td
					}
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
