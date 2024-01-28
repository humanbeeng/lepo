package golang

import (
	"go/ast"
	"go/token"
	"go/types"
	"log/slog"
	"strings"

	extract "github.com/humanbeeng/lepo/prototypes/analysis/extract"
)

type TypeVisitor struct {
	ast.Visitor
	Fset       *token.FileSet
	Info       *types.Info
	TypeDecls  map[string]extract.TypeDecl
	Implements map[string][]string
	Members    map[string]extract.Member
	Package    string
}

func (v *TypeVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch nd := node.(type) {

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

						stCode, err := extractCode(nd, v.Fset)
						if err != nil {
							// TODO: Handle errors gracefully
							panic(err)
						}
						impl, ok := v.Implements[stQName]
						if !ok {
							// Try with pointer type
							stQNameWPtr := "*" + stQName
							impl = v.Implements[stQNameWPtr]
						}

						td := extract.TypeDecl{
							Name:            tSpec.Name.Name,
							QName:           stQName,
							TypeQName:       tspecObj.Type().String(),
							Underlying:      tspecObj.Type().Underlying().String(),
							ImplementsQName: impl,
							Kind:            extract.Struct,
							Pos:             pos,
							End:             end,
							Filepath:        filepath,
							Code:            stCode,
							Doc: extract.Doc{
								Comment: nd.Doc.Text(),
								OfQName: stQName,
							},
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
					} else if id, ok := tSpec.Type.(*ast.Ident); ok {
						// Handle type aliases

						qname := tspecObj.Pkg().Path() + "." + tspecObj.Name()

						pos := v.Fset.Position(id.Pos()).Line
						end := v.Fset.Position(id.End()).Line
						filepath := v.Fset.Position(id.Pos()).Filename

						doc := extract.Doc{
							Comment: nd.Doc.Text(),
							OfQName: qname,
						}
						td := extract.TypeDecl{
							Name:       tspecObj.Name(),
							QName:      qname,
							TypeQName:  tspecObj.Type().String(),
							Underlying: tspecObj.Type().Underlying().String(),
							// TODO: Extract code
							Code:     "",
							Doc:      doc,
							Kind:     extract.Alias,
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
		return v
	}
}

func (v *TypeVisitor) handleFieldNode(field *ast.Field, parentQName string) error {
	if field == nil {
		return nil
	}

	for _, fieldName := range field.Names {
		fieldObj := v.Info.Defs[fieldName]
		fieldQName := parentQName + "." + fieldObj.Name()

		d := extract.Doc{
			Comment: field.Doc.Text() + field.Comment.Text(),
			OfQName: fieldQName,
			// TODO: Add doc type
		}

		if field.Tag != nil {
			d.Comment = d.Comment + field.Tag.Value
		}
		st, ok := field.Type.(*ast.StructType)
		if ok && (strings.HasPrefix(fieldObj.Type().String(), "struct")) {
			pos := v.Fset.Position(field.Pos())
			end := v.Fset.Position(field.End())

			var stCode string
			stCode, err := extractCode(st, v.Fset)
			if err != nil {
				return err
			}

			ftd := extract.TypeDecl{
				Name:       fieldObj.Name(),
				QName:      fieldQName,
				TypeQName:  fieldObj.Type().String(),
				Underlying: fieldObj.Type().Underlying().String(),
				Kind:       extract.Struct,
				Code:       stCode,
				Pos:        pos.Line,
				Doc:        d,
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
		m := extract.Member{
			Name:        fieldName.Name,
			QName:       fieldQName,
			TypeQName:   fieldObj.Type().String(),
			ParentQName: parentQName,
			Pos:         v.Fset.Position(field.Pos()).Line,
			End:         v.Fset.Position(field.End()).Line,
			Filepath:    v.Fset.Position(field.Pos()).Filename,
			// TODO: Extract member code
			Code: "",
			Doc:  d,
		}
		v.Members[fieldQName] = m
	}
	return nil
}
