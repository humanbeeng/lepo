package extract

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"log/slog"

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

		sv := &StructVisitor{
			Fset:      fset,
			Info:      pkg.TypesInfo,
			TypeDecls: make(map[string]TypeDecl),
			Members:   make(map[string]Member),
		}

		// iv := &InterfaceVisitor{
		// 	Fset:      fset,
		// 	Info:      pkg.TypesInfo,
		// 	TypeDecls: make(map[string]TypeDecl),
		// 	Members:   make(map[string]Member),
		// }

		// mv := &MethodVisitor{
		// 	Fset:    fset,
		// 	Info:    pkg.TypesInfo,
		// 	Methods: make(map[string]Function),
		// }

		if strings.Contains(pkg.PkgPath, pkgstr) {
			for _, syn := range pkg.Syntax {
				ast.Walk(sv, syn)
			}
		}

		// for _, m := range sv.TypeDecls {
		// 	pp.Println("TD", m.Name)
		// }

		for _, v := range sv.TypeDecls {
			slog.Info("Struct", "Name", v.Name)
			fmt.Println("------------")
		}

		for _, v := range sv.Members {
			slog.Info("Member", "Name", v.Name)
			fmt.Println("------------")
		}
	})

	return nil
}
