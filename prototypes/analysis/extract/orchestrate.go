package extract

import (
	"fmt"
	"go/ast"
	"go/token"
	"log/slog"
	"strings"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/refactor/satisfy"
)

func (g GoExtractor) Extract(pkgstr string) error {
	// TODO: Change or add directory path as well.

	// start := time.Now()

	slog.Info("Extraction requested for", "package", pkgstr)
	// same into cfg.Check method
	fset := token.NewFileSet()
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedDeps | packages.NeedSyntax |
			packages.NeedName | packages.NeedTypesInfo | packages.NeedImports,
		Fset: fset,
		Dir:  "/Users/apple/workspace/go/lepo/prototypes/go-testdata",
		// Dir: "/Users/apple/workspace/go/lepo/prototypes/analysis",
		// Dir: "/Users/apple/workspace/misc/dgraph",
	}

	// TODO: Take directory as input and get extract pkgstr using go mod file
	pkgs, err := packages.Load(cfg, pkgstr+"...")
	if err != nil {
		slog.Error("Unable to load", "package", pkgstr)
		return err
	}

	slog.Info("Found packages", "count", len(pkgs))

	implMap := make(map[string]string)

	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		// Process implementations of given package only.
		if !strings.Contains(pkg.PkgPath, pkgstr) {
			return
		}

		fi := satisfy.Finder{Result: make(map[satisfy.Constraint]bool)}
		fi.Find(pkg.TypesInfo, pkg.Syntax)

		// Transform Finder Result map to make it queryable
		for r := range fi.Result {
			implMap[r.RHS.String()] = r.LHS.String()
		}
	})

	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		// Process nodes of given package only.
		if !strings.Contains(pkg.PkgPath, pkgstr) {
			return
		}

		slog.Info("Analysing", "package", pkg.PkgPath)

		// iv := &InterfaceVisitor{
		// 	Fset:      fset,
		// 	Info:      pkg.TypesInfo,
		// 	TypeDecls: make(map[string]TypeDecl),
		// 	Members:   make(map[string]Member),
		// }

		nv := &NamedVisitor{
			Fset:  fset,
			Info:  pkg.TypesInfo,
			Named: make(map[string]Named),
		}

		// fv := &FileVisitor{
		// 	Imports: make([]Import, 0),
		// 	Package: pkg.PkgPath,
		// 	Fset:    fset,
		// 	Info:    pkg.TypesInfo,
		// }

		fv := &FunctionVisitor{
			Fset:      fset,
			Info:      pkg.TypesInfo,
			Functions: make(map[string]Function),
		}

		// tv := &TypeVisitor{
		// 	Fset:         fset,
		// 	Info:         pkg.TypesInfo,
		// 	TypeDecls:    make(map[string]TypeDecl),
		// 	Implementors: implMap,
		// 	Members:      make(map[string]Member),
		// }

		slog.Info("Files found in", "package", pkg.PkgPath, "count", len(pkg.Syntax))

		// For each file in package
		for _, file := range pkg.Syntax {
			// ast.Walk(fv, file)
			ast.Walk(fv, file)
			ast.Walk(nv, file)
		}
		// fmt.Println("Found", len(tv.TypeDecls), "types")

		// for _, m := range fv.Functions {
		// 	fmt.Println("Name:", m.Name)
		// 	for _, c := range m.Calls {
		// 		fmt.Println("Calls", c)
		// 	}
		// }

		for _, t := range nv.Named {
			fmt.Println("----------")
			fmt.Println("Name", t.Name)
			fmt.Println("QName", t.QName)
			fmt.Println("Type", t.TypeQName)
			fmt.Println("Underlying", t.Underlying)
			fmt.Println("Comment", t.Doc.Comment)
			fmt.Println("----------")
		}
	})

	return nil
}
