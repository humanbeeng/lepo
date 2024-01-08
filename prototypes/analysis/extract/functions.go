package extract

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

type FunctionVisitor struct {
	Functions map[string]Function
	Fset      *token.FileSet
	Info      *types.Info
}

func (v *FunctionVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch n := node.(type) {

	// TODO: Add FuncType which is suspected to be in interface
	case *ast.FuncDecl:
		{
			fnObj, ok := v.Info.Defs[n.Name]
			if !ok {
				return v
			}

			pos := v.Fset.Position(n.Pos()).Line
			end := v.Fset.Position(n.End()).Line
			filepath := v.Fset.Position(n.Pos()).Filename
			qname := fnObj.Pkg().Path() + "." + fnObj.Name()

			mCode, err := extractCode(n, v.Fset)
			if err != nil {
				// TODO: Better error handling
				panic(err)
			}

			if n.Recv != nil {
				for _, field := range n.Recv.List {
					if id, ok := field.Type.(*ast.Ident); ok {
						// Regular method
						stQName := fnObj.Pkg().Path() + "." + id.Name

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
						v.extractDoc(n, &f)

						v.Functions[qname] = f

					} else if se, ok := field.Type.(*ast.StarExpr); ok {
						// Pointer based method
						if id, ok := se.X.(*ast.Ident); ok {
							stQName := fnObj.Pkg().Path() + "." + id.Name

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
							v.extractDoc(n, &f)

							v.Functions[qname] = f

						}
					}
				}
			} else {
				// Just a regular function
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
				v.extractDoc(n, &f)

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
		return v
	}
}

func (v *FunctionVisitor) extractDoc(n *ast.FuncDecl, f *Function) {
	if n.Doc == nil {
		return
	}
	d := Doc{
		Comment: n.Doc.Text(),
		OfQName: f.QName,
	}
	if len(n.Doc.List) == 1 {
		d.Type = SingleLine
	} else if strings.HasPrefix(n.Doc.Text(), "/*") && strings.HasSuffix(n.Doc.Text(), "*/") {
		d.Type = Block
	} else {
		d.Type = MultiLine
	}

	f.Doc = d
}

func (v *FunctionVisitor) extractParamsAndReturns(n *ast.FuncDecl, f *Function) {
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