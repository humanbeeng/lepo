package golang

import (
	"go/ast"
	"go/token"
	"log/slog"
	"strings"
	"time"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/refactor/satisfy"

	"github.com/humanbeeng/lepo/prototypes/analysis/extract"
)

type GoExtractor struct {
	TypeDecls map[string]*extract.TypeDecl
	Members   map[string]*extract.Member
}

func NewGoExtractor() *GoExtractor {
	return &GoExtractor{
		TypeDecls: make(map[string]*extract.TypeDecl),
		Members:   make(map[string]*extract.Member),
	}
}

func (g *GoExtractor) Extract(pkgstr string, dir string) (extract.ExtractResult, error) {
	// TODO: Change or add directory path as well.
	start := time.Now()

	slog.Info("Extraction requested for", "package", pkgstr)
	// same into cfg.Check method
	fset := token.NewFileSet()
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedDeps | packages.NeedSyntax |
			packages.NeedName | packages.NeedTypesInfo | packages.NeedImports,
		Fset: fset,
		Dir:  dir,
	}

	// TODO: Take directory as input and get extract pkgstr using go mod file
	pkgs, err := packages.Load(cfg, pkgstr+"...")
	if err != nil {
		slog.Error("Unable to load", "package", pkgstr)
		return extract.ExtractResult{}, err
	}

	slog.Info("Packages found", "count", len(pkgs))

	implMap := make(map[string][]string)

	extractRes := extract.ExtractResult{
		TypeDecls:  make(map[string]extract.TypeDecl),
		Interfaces: make(map[string]extract.TypeDecl),
		NamedTypes: make(map[string]extract.Named),
		Members:    make(map[string]extract.Member),
		Functions:  make(map[string]extract.Function),
		Files:      make(map[string]extract.File),
	}

	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		// Process implementations of given package only.
		if !strings.Contains(pkg.PkgPath, pkgstr) {
			return
		}

		slog.Info("Constructing implementors map", "package", pkg.PkgPath)

		fi := satisfy.Finder{Result: make(map[satisfy.Constraint]bool)}
		fi.Find(pkg.TypesInfo, pkg.Syntax)

		// Transform Finder Result map to make it queryable
		for r := range fi.Result {
			implMap[r.RHS.String()] = append(implMap[r.RHS.String()], r.LHS.String())
		}
	})

	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		// Process nodes of given package only.
		if !strings.Contains(pkg.PkgPath, pkgstr) {
			return
		}

		slog.Info("Analysing", "package", pkg.PkgPath)
		slog.Info("Files found in", "package", pkg.PkgPath, "count", len(pkg.Syntax))

		tv := &TypeVisitor{
			Fset:       fset,
			Info:       pkg.TypesInfo,
			TypeDecls:  extractRes.TypeDecls,
			Implements: implMap,
			Members:    extractRes.Members,
			Package:    pkg.PkgPath,
		}

		nv := &NamedVisitor{
			Fset:  fset,
			Info:  pkg.TypesInfo,
			Named: extractRes.NamedTypes,
		}

		fv := &FunctionVisitor{
			Fset:      fset,
			Info:      pkg.TypesInfo,
			Functions: extractRes.Functions,
		}

		fiv := &FileVisitor{
			Package: pkg.PkgPath,
			Fset:    fset,
			Info:    pkg.TypesInfo,
			Files:   extractRes.Files,
		}

		iv := &InterfaceVisitor{
			Fset:       fset,
			Info:       pkg.TypesInfo,
			Interfaces: extractRes.Interfaces,
			Members:    extractRes.Members,
		}

		for _, file := range pkg.Syntax {
			slog.Info("Walking", "file", fset.Position(file.Pos()).Filename)
			ast.Walk(tv, file)
			ast.Walk(nv, file)
			ast.Walk(fv, file)
			ast.Walk(fiv, file)
			ast.Walk(iv, file)
		}
	})

	slog.Info("Extraction completed", "time_taken", time.Since(start).String())

	return extractRes, nil
}
