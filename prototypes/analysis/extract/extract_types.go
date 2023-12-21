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

	switch nd := node.(type) {

	case *ast.File:
		return v

	case *ast.Field:
		{
			for _, fieldName := range nd.Names {
				fieldObj := v.Info.Defs[fieldName]

				if st, ok := nd.Type.(*ast.StructType); ok {

					stQName := fieldObj.Pkg().Path() + "." + fieldName.Name
					pos := v.Fset.Position(st.Pos()).Line
					end := v.Fset.Position(st.End()).Line
					filepath := v.Fset.Position(st.Pos()).Filename

					// Extract code
					var stCode string
					var b []byte
					buf := bytes.NewBuffer(b)
					err := format.Node(buf, v.Fset, st)
					if err != nil {
						// TODO: Handle errors gracefully
						panic(err)
					}
					stCode = buf.String()

					td := TypeDecl{
						Name:       fieldName.Name,
						QName:      stQName,
						Type:       fieldObj.Type().String(),
						Underlying: fieldObj.Type().Underlying().String(),
						Kind:       Struct,
						Pos:        pos,
						End:        end,
						Filepath:   filepath,
						Code:       stCode,
					}

					v.TypeDecls[stQName] = td

					fields := st.Fields

					for _, field := range fields.List {
						// v.handleFieldNode(field, stQName)
						ast.Walk(v, field)
					}
				}
			}
			return v
		}

	case *ast.GenDecl:
		{
			for _, s := range nd.Specs {
				if tSpec, ok := s.(*ast.TypeSpec); ok {
					tspecObj := v.Info.Defs[tSpec.Name]

					if st, ok := tSpec.Type.(*ast.StructType); ok {
						stQName := tspecObj.Pkg().Path() + "." + tSpec.Name.Name
						pos := v.Fset.Position(st.Pos()).Line
						end := v.Fset.Position(st.End()).Line
						filepath := v.Fset.Position(st.Pos()).Filename
						var stCode string

						var b []byte
						buf := bytes.NewBuffer(b)
						err := format.Node(buf, v.Fset, nd)
						if err != nil {
							// TODO: Handle errors gracefully
							panic(err)
						}

						stCode = buf.String()

						td := TypeDecl{
							Name:       tSpec.Name.Name,
							QName:      stQName,
							Type:       tspecObj.Type().String(),
							Underlying: tspecObj.Type().Underlying().String(),
							Kind:       Struct,
							Pos:        pos,
							End:        end,
							Filepath:   filepath,
							Code:       stCode,
						}

						v.TypeDecls[stQName] = td

						fields := st.Fields
						for _, field := range fields.List {
							v.handleFieldNode(field, stQName)
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

func (v *StructVisitor) handleFieldNode(field *ast.Field, parentQName string) {
	for _, fieldName := range field.Names {
		fieldObj := v.Info.Defs[fieldName]
		fieldQName := parentQName + "." + fieldObj.Name()
		st, ok := field.Type.(*ast.StructType)
		if ok {
			if strings.HasPrefix(fieldObj.Type().String(), "struct") {
				pos := v.Fset.Position(field.Pos())
				end := v.Fset.Position(field.End())

				var stCode string

				var b []byte
				buf := bytes.NewBuffer(b)
				err := format.Node(buf, v.Fset, st)
				if err != nil {
					// TODO: Handle errors gracefully
					panic(err)
				}

				stCode = buf.String()
				ftd := TypeDecl{
					Name:       fieldObj.Name(),
					QName:      fieldQName,
					Type:       fieldObj.Type().String(),
					Underlying: fieldObj.Type().Underlying().String(),
					Kind:       Struct,
					Code:       stCode,
					Pos:        pos.Line,
					End:        end.Line,
					Filepath:   pos.Filename,
				}
				v.TypeDecls[fieldQName] = ftd
			}
			ast.Walk(v, field)
		}
		m := Member{
			Name:        fieldObj.Name(),
			QName:       fieldQName,
			TypeQName:   fieldObj.Type().String(),
			ParentQName: parentQName,
			Pos:         v.Fset.Position(field.Pos()).Line,
			End:         v.Fset.Position(field.End()).Line,
			Filepath:    v.Fset.Position(field.Pos()).Filename,
			Code:        "",
		}
		v.Members[fieldQName] = m
	}
}
