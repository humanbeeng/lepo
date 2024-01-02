package extract

import (
	"fmt"
	"go/ast"
	"go/token"
	"log/slog"
	"strings"

	"golang.org/x/tools/go/packages"
)

func (g *GoExtractor) Extract(pkgstr string) error {
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
		Dir:  "/Users/apple/workspace/go/lepo/prototypes/go-testdata",
	}

	// TODO: Take directory as input and get extract pkgstr using go mod file
	pkgs, err := packages.Load(cfg, pkgstr+"...")
	if err != nil {
		slog.Error("Unable to load", "package", pkgstr)
		return err
	}

	slog.Info("Found packages", "count", len(pkgs))
	// ivs := 0
	// tds := 0
	// fxs := 0

	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		// TODO : Move this inside if condition below
		// If this is your own package, process its structs.
		if !strings.Contains(pkg.PkgPath, pkgstr) {
			return
		}
		slog.Info("Analysing", "package", pkg.PkgPath)
		// tv := &TypeVisitor{
		// 	Fset:      fset,
		// 	Info:      pkg.TypesInfo,
		// 	TypeDecls: make(map[string]TypeDecl),
		// 	Members:   make(map[string]Member),
		// }
		//
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
		// 	Imports: make([]string, 0),
		// 	Package: pkg.PkgPath,
		// 	Fset:    fset,
		// 	Info:    pkg.TypesInfo,
		// }

		mv := &MethodVisitor{
			Fset:      fset,
			Info:      pkg.TypesInfo,
			Functions: make(map[string]Function),
		}
		// var wg *sync.WaitGroup

		// For each file in package
		for _, syn := range pkg.Syntax {
			ast.Walk(mv, syn)
			// ast.Walk(sv, syn)
			ast.Walk(iv, syn)
		}
		// fmt.Println("Found", len(sv.TypeDecls), "types")

		// for _, c := range tv.TypeDecls {
		// 	fmt.Println(c.Name)
		// 	fmt.Println(c.QName)
		// 	fmt.Println(c.TypeQName)
		// 	fmt.Println(c.Underlying)
		// 	fmt.Println(c.Kind)
		// 	fmt.Println("-------")
		// }

		for _, m := range mv.Functions {
			fmt.Println("Function Name:", m.QName)
		}

		// tds += len(sv.TypeDecls)
		// ivs += len(iv.TypeDecls)
		// fxs += len(mv.Functions)
	})

	// slog.Info("TypeDecls", "count", tds)
	// slog.Info("Interfaces", "count", ivs)
	// slog.Info("Functions", "count", fxs)
	// slog.Info("Extraction completed", "time", time.Since(start).Seconds())
	return nil
}
