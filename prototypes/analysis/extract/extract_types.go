package extract

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"strings"
)

type StructVisitor struct {
	ast.Visitor
	Fset      *token.FileSet
	Info      *types.Info
	TypeDecls map[string]TypeDecl
	Members   map[string]Member
	Files     map[string][]byte
}

func (v *StructVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch n := node.(type) {

	case *ast.File,
		*ast.FieldList,
		*ast.Ident:
		return v

	case *ast.GenDecl:
		{
			for _, s := range n.Specs {
				if ts, ok := s.(*ast.TypeSpec); ok {
					tsObj := v.Info.Defs[ts.Name]

					if st, ok := ts.Type.(*ast.StructType); ok {
						stQName := tsObj.Pkg().Path() + "." + ts.Name.Name
						pos := v.Fset.Position(st.Pos()).Line
						end := v.Fset.Position(st.End()).Line
						filepath := v.Fset.Position(st.Pos()).Filename
						var stCode string

						var b []byte
						buf := bytes.NewBuffer(b)
						err := format.Node(buf, v.Fset, n)
						if err != nil {
							panic(err)
						}

						stCode = buf.String()

						td := TypeDecl{
							Name:       ts.Name.Name,
							QName:      stQName,
							Type:       tsObj.Type().String(),
							Underlying: tsObj.Type().Underlying().String(),
							Kind:       Struct,
							Pos:        pos,
							End:        end,
							Filepath:   filepath,
							Code:       stCode,
						}

						v.TypeDecls[stQName] = td

						fields := st.Fields

						for _, f := range fields.List {
							for _, n := range f.Names {
								fobj := v.Info.Defs[n]
								fQName := fobj.Pkg().Path() + "." + fobj.Name()
								_, ok := fobj.Type().Underlying().(*types.Struct)
								if ok {
									// Store only structs are are defined in project
									if strings.HasPrefix(fobj.Type().String(), "struct") {
										pos := v.Fset.Position(f.Pos())
										end := v.Fset.Position(f.End())

										ftd := TypeDecl{
											Name:       fobj.Name(),
											QName:      fQName,
											Type:       fobj.Type().String(),
											Underlying: fobj.Type().Underlying().String(),
											Kind:       Struct,
											Pos:        pos.Line,
											End:        end.Line,
											Filepath:   pos.Filename,
										}
										v.TypeDecls[fQName] = ftd
									}
								}
								m := Member{
									Name:        fobj.Name(),
									QName:       fQName,
									TypeQName:   fobj.Type().String(),
									ParentQName: stQName,
									Pos:         v.Fset.Position(f.Pos()).Line,
									End:         v.Fset.Position(f.End()).Line,
									Filepath:    v.Fset.Position(f.Pos()).Filename,
									Code:        "",
								}
								v.Members[fQName] = m
							}
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
