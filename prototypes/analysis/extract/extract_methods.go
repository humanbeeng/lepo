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

	// TODO: Add FuncType which is suspected to be in interface
	case *ast.FuncDecl:
		{
			// Add nil check for obj
			fnObj := v.Info.Defs[n.Name]

			pos := v.Fset.Position(n.Pos()).Line
			end := v.Fset.Position(n.End()).Line
			filepath := v.Fset.Position(n.Pos()).Filename

			var mCode string
			var b []byte

			buf := bytes.NewBuffer(b)
			err := format.Node(buf, v.Fset, n)
			if err != nil {
				// TODO: Better error handling
				panic(err)
			}

			mCode = buf.String()
			if n.Recv != nil {
				for _, field := range n.Recv.List {
					if id, ok := field.Type.(*ast.Ident); ok {
						// Regular method
						stQName := fnObj.Pkg().Path() + "." + id.Name
						qname := fnObj.Pkg().Path() + "." + fnObj.Name()

						f := Function{
							Name:        fnObj.Name(),
							QName:       qname,
							ParentQName: stQName,
							Pos:         pos,
							End:         end,
							Filepath:    filepath,
							Code:        mCode,
						}

						v.Methods[qname] = f
					} else if se, ok := field.Type.(*ast.StarExpr); ok {
						// Pointer based method
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
								Code:        mCode,
							}
							v.Methods[qname] = f
						}
					}
				}
			} else {
				// Just a regular function
				qname := fnObj.Pkg().Path() + "." + fnObj.Name()
				f := Function{
					Name:        fnObj.Name(),
					QName:       qname,
					ParentQName: "",
					Pos:         pos,
					End:         end,
					Filepath:    filepath,
					Code:        mCode,
				}

				v.Methods[qname] = f
			}

			// ast.Print(v.Fset, n.Body)
			bv := &BodyVisitor{
				Fset: v.Fset,
				Info: v.Info,
			}

			ast.Walk(bv, n.Body)
			return v
		}

	default:
		return nil
	}
}

func (v *MethodVisitor) handleCallExpr(ce *ast.CallExpr) {
	if id, ok := ce.Fun.(*ast.Ident); ok {
		ceObj := v.Info.Uses[id]
		if ceObj != nil {
			// qname := ceObj.Pkg().Path() + "." + ceObj.Name()
			println("Calling", ceObj.Name())
		}
	} else if se, ok := ce.Fun.(*ast.SelectorExpr); ok {
		seObj := v.Info.Uses[se.Sel]
		println("Calling", seObj.Name())
	}
}
