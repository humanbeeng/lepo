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

							println("Method:", f.QName)
							v.Methods[qname] = f
						}
					}
				}
			} else {
				// Not a method, just a regular function
				qname := fnObj.Pkg().Path() + "." + fnObj.Name()
				f := Function{
					Name:        fnObj.Name(),
					QName:       qname,
					ParentQName: "",
					Pos:         pos,
					End:         end,
					Filepath:    filepath,
					Code:        methCode,
				}

				v.Methods[qname] = f
				// println("Method", qname)
			}

			// ast.Print(v.Fset, n.Body)
			bv := &BodyVisitor{
				Fset: v.Fset,
				Info: v.Info,
			}

			ast.Walk(bv, n.Body)
			// for _, stmt := range n.Body.List {
			// 	if as, ok := stmt.(*ast.AssignStmt); ok {
			// 		for _, e := range as.Rhs {
			// 			if ce, ok := e.(*ast.CallExpr); ok {
			// 				v.handleCallExpr(ce)
			// 			}
			// 		}
			// 	} else if es, ok := stmt.(*ast.ExprStmt); ok {
			// 		if ce, ok := es.X.(*ast.CallExpr); ok {
			// 			// v.handleCallExpr(ce)
			// 		}
			// 	}
			// }
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
