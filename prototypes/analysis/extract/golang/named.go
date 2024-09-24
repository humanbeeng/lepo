package golang

import (
	"go/ast"
	"go/token"
	"go/types"
	"log/slog"

	"github.com/humanbeeng/lepo/prototypes/analysis/extract"
)

type NamedVisitor struct {
	ast.Visitor
	Named map[string]extract.Named
	Fset  *token.FileSet
	Info  *types.Info
}

func (v *NamedVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch nd := node.(type) {
	case *ast.GenDecl:
		{
			for _, spec := range nd.Specs {
				// TODO: Refactor this nested ifs
				ts, tsOk := spec.(*ast.TypeSpec)
				if tsOk {
					_, ftOk := ts.Type.(*ast.FuncType)
					if ftOk {
						pos := v.Fset.Position(ts.Pos()).Line
						end := v.Fset.Position(ts.End()).Line
						ftObj, objOk := v.Info.Defs[ts.Name]
						if objOk {
							namespace := extract.Namespace{Name: ftObj.Pkg().Path()}
							ftQName := namespace.Name + "." + ftObj.Name()
							codeStr, err := extractCode(ts, v.Fset)
							if err != nil {
								slog.Error("Unable to extract code", "qname", ftQName)
							}
							funcType := extract.Named{
								Name:       ftObj.Name(),
								QName:      ftQName,
								Namespace:  namespace,
								TypeQName:  ftObj.Type().String(),
								Underlying: ftObj.Type().Underlying().String(),
								Pos:        pos,
								End:        end,
								Code:       codeStr,
								Doc: extract.Doc{
									Comment: ts.Doc.Text() + ts.Comment.Text(),
									OfQName: ftQName,
									// TODO: Revisit
									Type: extract.SingleLine,
								},
							}
							v.Named[ftQName] = funcType
						}

					}
				}

				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				for _, name := range vs.Names {
					vsObj := v.Info.Defs[name]

					pos := v.Fset.Position(vs.Pos()).Line
					end := v.Fset.Position(vs.End()).Line
					filepath := v.Fset.Position(vs.Pos()).Filename

					code, _ := extractCode(node, v.Fset)
					namespace := extract.Namespace{Name: vsObj.Pkg().Path()}
					qname := namespace.Name + "." + name.Name

					named := extract.Named{
						Name:       name.Name,
						QName:      qname,
						TypeQName:  vsObj.Type().String(),
						Underlying: vsObj.Type().Underlying().String(),
						Code:       code,
						Doc:        extract.Doc{Comment: vs.Doc.Text() + vs.Comment.Text(), OfQName: qname},
						Pos:        pos,
						End:        end,
						Filepath:   filepath,
						Namespace:  namespace,
					}
					v.Named[qname] = named
				}

			}
		}
	default:
		{
			return v
		}
	}
	return nil
}
