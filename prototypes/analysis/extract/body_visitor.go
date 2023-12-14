package extract

import (
	"go/ast"
	"go/token"
	"go/types"
)

type BodyVisitor struct {
	ast.Visitor
	Fset *token.FileSet
	Info *types.Info
}

func (v *BodyVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *ast.BlockStmt:
		return v

	case *ast.AssignStmt:
		println("Assign statement called")
		{
			for _, e := range n.Rhs {
				if ce, ok := e.(*ast.CallExpr); ok {
					v.handleCallExpr(ce)
				}
			}
		}

	case *ast.ExprStmt:
		{
			if ce, ok := n.X.(*ast.CallExpr); ok {
				v.handleCallExpr(ce)
			}
		}

	case *ast.IfStmt:
		{
		}

	case *ast.ReturnStmt:
		{
		}

	case *ast.ForStmt:
		{
		}

	default:
		return nil
	}
	return nil
}

func (v *BodyVisitor) handleCallExpr(ce *ast.CallExpr) {
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
