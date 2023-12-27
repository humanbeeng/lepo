package extract

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
)

type MethodVisitor struct {
	Functions map[string]Function
	Fset      *token.FileSet
	Info      *types.Info
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
			var qname string

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
						qname = fnObj.Pkg().Path() + "." + fnObj.Name()

						f := Function{
							Name:        fnObj.Name(),
							QName:       qname,
							ParentQName: stQName,
							Pos:         pos,
							End:         end,
							Filepath:    filepath,
							Code:        mCode,
						}

						v.extractParamsAndReturns(n, &f)
						v.Functions[qname] = f
					} else if se, ok := field.Type.(*ast.StarExpr); ok {
						// Pointer based method
						if id, ok := se.X.(*ast.Ident); ok {
							stQName := fnObj.Pkg().Path() + "." + id.Name
							qname = fnObj.Pkg().Path() + "." + fnObj.Name()

							f := Function{
								Name:        fnObj.Name(),
								QName:       qname,
								ParentQName: stQName,
								Pos:         pos,
								End:         end,
								Filepath:    filepath,
								Code:        mCode,
							}
							v.extractParamsAndReturns(n, &f)
							v.Functions[qname] = f
						}
					}
				}
			} else {
				// Just a regular function
				qname = fnObj.Pkg().Path() + "." + fnObj.Name()
				f := Function{
					Name:        fnObj.Name(),
					QName:       qname,
					ParentQName: "",
					Pos:         pos,
					End:         end,
					Filepath:    filepath,
					Code:        mCode,
				}

				v.extractParamsAndReturns(n, &f)
				v.Functions[qname] = f
			}

			bv := &BodyVisitor{
				CallerQName: qname,
				Fset:        v.Fset,
				Info:        v.Info,
			}
			if n.Body != nil {
				ast.Walk(bv, n.Body)
			}
			return v
		}

	default:
		return nil
	}
}

func (v *MethodVisitor) extractParamsAndReturns(n *ast.FuncDecl, f *Function) {
	if n.Type.Params == nil {
		return
	}
	params := n.Type.Params.List
	for _, p := range params {
		for _, name := range p.Names {
			pObj := v.Info.Defs[name]
			f.ParamQNames = append(f.ParamQNames, pObj.Type().String())
		}
	}

	if n.Type.Results == nil {
		return
	}
	results := n.Type.Results.List
	for _, r := range results {
		a, ok := v.Info.Types[r.Type]
		if !ok {
			continue
		}
		f.ReturnQNames = append(f.ReturnQNames, a.Type.String())
	}
}
