package extract

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"log/slog"
	"strings"
)

type TypeVisitor struct {
	ast.Visitor
	Fset      *token.FileSet
	Info      *types.Info
	TypeDecls map[string]TypeDecl
	Members   map[string]Member
}

func (v *TypeVisitor) Visit(node ast.Node) ast.Visitor {
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

						stCode, err := code(nd, v.Fset)
						if err != nil {
							// TODO: Handle errors gracefully
							panic(err)
						}

						td := TypeDecl{
							Name:       tSpec.Name.Name,
							QName:      stQName,
							TypeQName:  tspecObj.Type().String(),
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
							err := v.handleFieldNode(field, stQName)
							// TODO: Revisit on how to handle errors
							if err != nil {
								slog.Error("Unable to visit field", err)
								fmt.Print("")
							}
						}
					} else if id, ok := tSpec.Type.(*ast.Ident); ok {
						// Handle type aliases

						qname := tspecObj.Pkg().Path() + "." + tspecObj.Name()

						pos := v.Fset.Position(id.Pos()).Line
						end := v.Fset.Position(id.End()).Line
						filepath := v.Fset.Position(id.Pos()).Filename
						td := TypeDecl{
							Name:       tspecObj.Name(),
							QName:      qname,
							TypeQName:  tspecObj.Type().String(),
							Underlying: tspecObj.Type().Underlying().String(),
							// TODO: Extract code
							Code:     "",
							Kind:     Alias,
							Pos:      pos,
							End:      end,
							Filepath: filepath,
						}
						v.TypeDecls[qname] = td
					}
				}
			}
			return v
		}

	default:
		return nil
	}
}

func (v *TypeVisitor) handleFieldNode(field *ast.Field, parentQName string) error {
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
				TypeQName:  fieldObj.Type().String(),
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
