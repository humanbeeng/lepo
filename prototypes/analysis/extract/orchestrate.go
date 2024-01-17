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

	// orchestrate extract
	slog.Info("Extraction requested for", "package", pkgstr)
	// same into cfg.Check method
	fset := token.NewFileSet()
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedDeps | packages.NeedSyntax |
			packages.NeedName | packages.NeedTypesInfo | packages.NeedImports,
		Fset: fset,
		// Dir:  "/Users/apple/workspace/go/lepo/prototypes/go-testdata",
		Dir: "/Users/apple/workspace/misc/dgraph",
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
	fmt.Println("Len of impl", len(implMap))

	// for k, v := range implMap {
	// 	fmt.Println(k, "implements", v)
	// 	fmt.Println("-------")
	// }

	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		// If this is your own package, process its structs.
		if !strings.Contains(pkg.PkgPath, pkgstr) {
			return
		}

		slog.Info("Analysing", "package", pkg.PkgPath)

		iv := &InterfaceVisitor{
			Fset:      fset,
			Info:      pkg.TypesInfo,
			TypeDecls: make(map[string]TypeDecl),
			Members:   make(map[string]Member),
		}

		// cv := &ConstVisitor{
		// 	Fset:      fset,
		// 	Info:      pkg.TypesInfo,
		// 	Constants: make(map[string]Constant),
		// }

		// fv := &FileVisitor{
		// 	Imports: make([]Import, 0),
		// 	Package: pkg.PkgPath,
		// 	Fset:    fset,
		// 	Info:    pkg.TypesInfo,
		// }

		// fv := &FunctionVisitor{
		// 	Fset:      fset,
		// 	Info:      pkg.TypesInfo,
		// 	Functions: make(map[string]Function),
		// }

		// slog.Info("Files found", "count", len(pkg.Syntax))

		tv := &TypeVisitor{
			Fset:         fset,
			Info:         pkg.TypesInfo,
			TypeDecls:    make(map[string]TypeDecl),
			Implementors: implMap,
			Members:      make(map[string]Member),
		}

		// For each file in package

		for _, file := range pkg.Syntax {
			ast.Walk(tv, file)
			ast.Walk(iv, file)
		}
		// fmt.Println("Found", len(tv.TypeDecls), "types")

		for _, c := range tv.TypeDecls {
			if c.ImplementsQName == "" {
				continue
			}
			fmt.Printf("-----\n\n")
			fmt.Println("Name", c.Name)
			fmt.Println("Implements", c.ImplementsQName)
			fmt.Printf("-----\n\n")
		}

		// for _, m := range iv.TypeDecls {
		// 	fmt.Println("Name:", m.Name)
		// 	fmt.Println("TypeQName:", m.TypeQName)
		// }
		// for _, m := range fv.Imports {
		// 	fmt.Printf("%+v\n", m)
		// }
	})

	return nil
}
