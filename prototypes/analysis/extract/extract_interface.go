package extract

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"

	"github.com/k0kubun/pp"
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
	case *ast.File,
		*ast.Ident,
		*ast.FieldList:
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

						var infCode string
						var b []byte

						buf := bytes.NewBuffer(b)
						err := format.Node(buf, v.Fset, n)
						if err != nil {
							panic(err)
						}

						infCode = buf.String()

						td := TypeDecl{
							Name:       ts.Name.Name,
							QName:      infQname,
							Type:       tsObj.Type().String(),
							Underlying: tsObj.Type().Underlying().String(),
							Kind:       Interface,
							Pos:        pos,
							End:        end,
							Filepath:   filepath,
							Code:       infCode,
						}
						pp.Println("Interface", td)
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
