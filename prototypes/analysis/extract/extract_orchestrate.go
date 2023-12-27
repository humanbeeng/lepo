package extract

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"time"

	"log/slog"

	"golang.org/x/tools/go/packages"
)

func (g *GoExtractor) Extract(pkgstr string) error {
	// TODO: Change or add directory path as well.
	start := time.Now()
	// orchestrate extract
	fmt.Println("Extraction requested for", pkgstr)
	// same into cfg.Check method
	fset := token.NewFileSet()
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedDeps | packages.NeedSyntax |
			packages.NeedName | packages.NeedTypesInfo | packages.NeedImports,
		Fset: fset,
		Dir:  "/home/humanbeeng/projects/read-only/dgraph",
	}

	// TODO: Take directory as input and get extract pkgstr using go mod file
	pkgs, err := packages.Load(cfg, pkgstr+"...")
	if err != nil {
		fmt.Println("Unable to load package")
		return err
	}

	fmt.Println("Found", len(pkgs), "packages")
	ivs := 0
	tds := 0
	fxs := 0

	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		// If this is your own package, process its structs.
		// TODO : Move this inside if condition below

		if strings.Contains(pkg.PkgPath, pkgstr) {
			slog.Info("Analysing", "package", pkg.PkgPath)
			sv := &StructVisitor{
				Fset:      fset,
				Info:      pkg.TypesInfo,
				TypeDecls: make(map[string]TypeDecl),
				Members:   make(map[string]Member),
			}

			iv := &InterfaceVisitor{
				Fset:      fset,
				Info:      pkg.TypesInfo,
				TypeDecls: make(map[string]TypeDecl),
				Members:   make(map[string]Member),
			}

			mv := &MethodVisitor{
				Fset:      fset,
				Info:      pkg.TypesInfo,
				Functions: make(map[string]Function),
			}
			// var wg *sync.WaitGroup

			// For each file in package
			for _, syn := range pkg.Syntax {
				ast.Walk(mv, syn)
				ast.Walk(sv, syn)
				ast.Walk(iv, syn)
			}
			tds += len(sv.TypeDecls)
			ivs += len(iv.TypeDecls)
			fxs += len(mv.Functions)
		}
	})
	slog.Info("TypeDecls", "count", tds)
	slog.Info("Interfaces", "count", ivs)
	slog.Info("Functions", "count", fxs)
	slog.Info("Extraction completed", "time", time.Since(start).Seconds())
	return nil
}
