package extract

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"

	"strings"

	"log/slog"

)

type StructVisitor struct {
	ast.Visitor
	Fset      *token.FileSet
	Info      *types.Info
	TypeDecls map[string]TypeDef
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
						var code string

						var b []byte
						buf := bytes.NewBuffer(b)
						err := format.Node(buf, v.Fset, nd)
						if err != nil {

							// TODO: Handle errors gracefully

							panic(err)
						}

						code = buf.String()
						node := Node{
							Name:          ts.Name.Name,
							QualifiedName: qualified,
							Pos:           pos,
							End:           end,
							File:          filepath,
							Code:          code,
						}


						td := TypeDecl{
							Name:       tSpec.Name.Name,
							QName:      stQName,
							Type:       tspecObj.Type().String(),
							Underlying: tspecObj.Type().Underlying().String(),

							Kind:       Struct,
						}

						v.TypeDecls[stQName] = td

						fields := st.Fields
						for _, field := range fields.List {
							err := v.handleFieldNode(field, stQName)
							// TODO: Revisit on how to handle errors
							if err != nil {
								slog.Error("Unable to visit field", err)
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

func (v *StructVisitor) handleFieldNode(field *ast.Field, parentQName string) error {
	if field == nil {
		return nil
	}

	for _, fieldName := range field.Names {
		fieldObj := v.Info.Defs[fieldName]
		fieldQName := parentQName + "." + fieldObj.Name()

		st, ok := field.Type.(*ast.StructType)
		if ok && (strings.HasPrefix(fieldObj.Type().String(), "struct")) {
			pos := v.Fset.Position(field.Pos())
			end := v.Fset.Position(field.End())

			var stCode string

			var b []byte
			buf := bytes.NewBuffer(b)
			err := format.Node(buf, v.Fset, st)
			if err != nil {
				return err
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

			fields := st.Fields
			for _, stf := range fields.List {
				err := v.handleFieldNode(stf, fieldQName)
				if err != nil {
					return err
				}
			}

		}
		m := Member{
			Name:        fieldName.Name,
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
	return nil
}
