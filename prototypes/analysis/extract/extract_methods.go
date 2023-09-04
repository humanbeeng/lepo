package extract

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
)

type MethodVisitor struct {
	Methods map[string]Function
	Fset    *token.FileSet
	Info    *types.Info
}

func (v *MethodVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch n := node.(type) {

	case *ast.File:
		return v

	case *ast.FuncDecl:
		{
			// Check if it's a method declaration with a receiver
			// Add nil check for obj
			fnObj := v.Info.Defs[n.Name]

			pos := v.Fset.Position(n.Pos()).Line
			end := v.Fset.Position(n.End()).Line
			filepath := v.Fset.Position(n.Pos()).Filename

			var methCode string
			var b []byte

			buf := bytes.NewBuffer(b)
			err := format.Node(buf, v.Fset, n)
			if err != nil {
				// Remove this
				panic(err)
			}

			methCode = buf.String()
			if n.Recv != nil {
				for _, field := range n.Recv.List {
					if id, ok := field.Type.(*ast.Ident); ok {
						stQName := fnObj.Pkg().Path() + "." + id.Name
						qname := fnObj.Pkg().Path() + "." + fnObj.Name()

						f := Function{
							Name:        fnObj.Name(),
							QName:       qname,
							ParentQName: stQName,
							Pos:         pos,
							End:         end,
							Filepath:    filepath,
							Code:        methCode,
						}

						v.Methods[qname] = f
					} else if se, ok := field.Type.(*ast.StarExpr); ok {
						if id, ok := se.X.(*ast.Ident); ok {
							stQName := fnObj.Pkg().Path() + "." + id.Name
							qname := fnObj.Pkg().Path() + "." + fnObj.Name()

							f := Function{
								Name:        fnObj.Name(),
								QName:       qname,
								ParentQName: stQName,
								Pos:         pos,
								End:         end,
								Filepath:    filepath,
								Code:        methCode,
							}

							v.Methods[qname] = f
						}
					}
				}
			}
			return v
		}

	default:
		return nil
	}
}
