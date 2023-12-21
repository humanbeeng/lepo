package extract

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

var stdout io.Writer = os.Stdout

func (g *GoExtractor) Extract(pkgstr string) error {
	// orchestrate extract
	fmt.Println("Extraction requested for", pkgstr)

	// same into cfg.Check method
	fset := token.NewFileSet()
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedDeps | packages.NeedSyntax |
			packages.NeedName | packages.NeedTypesInfo | packages.NeedImports,
		Fset: fset,
		Dir:  "/Users/apple/projects/lepo/prototypes/analysis",
	}

	// TODO: Take directory as input and get extract pkgstr using go mod file
	pkgs, err := packages.Load(cfg, pkgstr+"...")
	if err != nil {
		fmt.Println("Unable to load package")
		return err
	}
	// prog, _ := ssautil.AllPackages(pkgs, ssa.PrintPackages|ssa.PrintFunctions)
	// prog.Build()
	// cg := static.CallGraph(prog)
	// var before, after string
	// format := "digraph"
	//
	// // Pre-canned formats.
	// switch format {
	// case "digraph":
	// 	format = `{{printf "%q %q" .Caller .Callee}}`
	//
	// case "graphviz":
	// 	before = "digraph callgraph {\n"
	// 	after = "}\n"
	// 	format = `  {{printf "%q" .Caller}} -> {{printf "%q" .Callee}}`
	// }
	//
	// tmpl, err := template.New("-format").Parse(format)
	// if err != nil {
	// 	return fmt.Errorf("invalid -format template: %v", err)
	// }

	// Allocate these once, outside the traversal.
	// var buf bytes.Buffer
	// data := callgraph.Edge{}
	// if err := callgraph.GraphVisitEdges(cg, func(edge *callgraph.Edge) error {
	// 	data.Caller = edge.Caller
	// 	data.Callee = edge.Callee
	// 	buf.Reset()
	// 	if err := tmpl.Execute(&buf, &data); err != nil {
	// 		return err
	// 	}
	// 	stdout.Write(buf.Bytes())
	// 	if len := buf.Len(); len == 0 || buf.Bytes()[len-1] != '\n' {
	// 		fmt.Fprintln(stdout)
	// 	}
	// 	return nil
	// }); err != nil {
	// 	return err
	// }
	// fmt.Fprint(stdout, after)
	// Create a call graph
	// Print the call graph

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

		// for _, m := range mv.Methods {
		// 	pp.Println("Method", m)
		// }

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

// mainPackages returns the main packages to analyze.
// Each resulting package is named "main" and has a main function.
// func mainPackages(pkgs []*ssa.Package) ([]*ssa.Package, error) {
// 	var mains []*ssa.Package
// 	for _, p := range pkgs {
// 		if p != nil && p.Pkg.Name() == "main" && p.Func("main") != nil {
// 			mains = append(mains, p)
// 		}
// 	}
// 	if len(mains) == 0 {
// 		return nil, fmt.Errorf("no main packages")
// 	}
// 	return mains, nil
// }
