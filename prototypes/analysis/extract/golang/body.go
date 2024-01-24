package golang

import (
	"go/ast"
	"go/token"
	"go/types"
)

type BodyVisitor struct {
	ast.Visitor
	CallerQName string
	Calls       []string
	Fset        *token.FileSet
	Info        *types.Info
}

func (v *BodyVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *ast.CallExpr:
		{
			v.handleCallExpr(n)
			return v
		}
	default:
		return v
	}
}

func (v *BodyVisitor) handleCallExpr(ce *ast.CallExpr) {
	if id, ok := ce.Fun.(*ast.Ident); ok {
		ceObj := v.Info.Uses[id]
		if ceObj != nil {
			callee := ceObj.Pkg().Path() + "." + ceObj.Name()
			v.Calls = append(v.Calls, callee)
		}
	} else if se, ok := ce.Fun.(*ast.SelectorExpr); ok {
		seObj := v.Info.Uses[se.Sel]
		if seObj != nil {
			callee := seObj.Pkg().Path() + "." + seObj.Name()
			v.Calls = append(v.Calls, callee)
		}
	}
}
