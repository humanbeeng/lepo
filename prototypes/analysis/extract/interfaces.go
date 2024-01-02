package extract

import (
	"go/ast"
	"go/token"
	"go/types"
)

type InterfaceVisitor struct {
	ast.Visitor
	Fset      *token.FileSet
	Info      *types.Info
	TypeDecls map[string]TypeDecl
	Members   map[string]Member
	Files     map[string][]byte
}

func (v *InterfaceVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *ast.File:
		return v

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

						infCode, err := code(n, v.Fset)
						if err != nil {
							panic(err)
						}

						td := TypeDecl{
							Name:       ts.Name.Name,
							QName:      infQname,
							TypeQName:  tsObj.Type().String(),
							Underlying: tsObj.Type().Underlying().String(),
							Kind:       Interface,
							Pos:        pos,
							End:        end,
							Filepath:   filepath,
							Code:       infCode,
						}
						v.TypeDecls[infQname] = td
					}
				}
			}
			return v
		}
	default:
		{
			return nil
		}
	}
}
