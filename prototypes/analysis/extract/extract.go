package extract

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/k0kubun/pp"
	"golang.org/x/tools/go/packages"
)

func (g *GoExtractor) Extract(pkgstr string) error {
	// orchestrate extract
	fmt.Println("Extraction requested for", pkgstr)

	// same into cfg.Check method
	fset := token.NewFileSet()
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedDeps | packages.NeedSyntax |
			packages.NeedName | packages.NeedTypesInfo | packages.NeedImports,
		Fset: fset,
		Dir:  "/home/humanbeeng/projects/lepo/prototypes/go-testdata",
	}

	// TODO: Take directory as input and get extract pkgstr using go mod file
	pkgs, err := packages.Load(cfg, pkgstr+"...")
	if err != nil {
		fmt.Println("Unable to load package")
		return err
	}

	fmt.Println("Found", len(pkgs), "packages")

	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		// If this is your own package, process its structs.
		// TODO : Move this inside if condition below

		// sv := &StructVisitor{
		// 	Fset:      fset,
		// 	Info:      pkg.TypesInfo,
		// 	TypeDecls: make(map[string]TypeDecl),
		// 	Members:   make(map[string]Member),
		// }
		// iv := &InterfaceVisitor{
		// 	Fset:      fset,
		// 	Info:      pkg.TypesInfo,
		// 	TypeDecls: make(map[string]TypeDecl),
		// 	Members:   make(map[string]Member),
		// }

		mv := &MethodVisitor{
			Fset:    fset,
			Info:    pkg.TypesInfo,
			Methods: make(map[string]Function),
		}

		if strings.Contains(pkg.PkgPath, pkgstr) {
			for _, syn := range pkg.Syntax {
				ast.Walk(mv, syn)
			}
		}

		for _, m := range mv.Methods {
			pp.Println("Method", m)
		}

		// for _, v := range sv.TypeDecls {
		// 	pp.Println("Struct:\n", v)
		// }
		//
		// for _, v := range sv.Members {
		// 	pp.Println("Member:", v)
		// }
	})

	return nil
}

// func (g *GoExtractor) ExtractMembers(st *types.Struct, parentQName string, fset *token.FileSet) {
// 	for i := 0; i < st.NumFields(); i++ {
// 		sf := st.Field(i)
// 		fieldName := sf.Name()
// 		fieldType := sf.Type().String()
// 		fieldQName := sf.Pkg().Path() + "." + sf.Name()
// 		if _, ok := g.TypeDecls[fieldQName]; ok {
// 			fmt.Println("Found already")
// 			break
// 		}
//
// 		m := &Member{
// 			Node: Node{
// 				Code: sf.String(),
// 				Pos:  fset.Position(sf.Pos()).Line,
// 				Name: fieldName,
// 			},
// 			ParentQName: parentQName,
// 			QualifiedName:       fieldQName,
// 			TypeQualifiedName:   fieldType,
// 		}
// 		g.Members[fieldQName] = m
// 	}
// }

// Dead code

// for _, obj := range pkg.TypesInfo.Defs {
// 	if obj == nil {
// 		continue
// 	}
// 	switch typ := obj.Type().(type) {
// 	case *types.Named:
// 		{
// 			if st, ok := typ.Underlying().(*types.Struct); ok {
// 				g.ExtractTypes(st, obj, fset)
// 				g.ExtractMembers(st, obj.Name(), fset)
// 			}
// 			break
// 		}
// 		// else if iface, ok := t.Underlying().(*types.Interface); ok {
// 		// 	// Case 2: interface
// 		// } else if sig, ok := t.Underlying().(*types.Signature); ok {
// 		// 	// Case 3: signature
// 		// }
//
// 	case *types.Struct:
// 		{
// 			g.ExtractTypes(typ, obj, fset)
// 			break
// 		}
// 	}
// }
