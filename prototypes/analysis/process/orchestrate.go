package process

import (
	"fmt"
	"log/slog"

	"github.com/humanbeeng/lepo/prototypes/analysis/extract"
)

// TODO: Refactor tf outta error handling and logging
func Orchestrate(e extract.Extractor) {
	slog.Info("Begin orchestration")
	// Step 1: Extract
	extractRes, err := e.Extract(
		"github.com/ardanlabs/service",
		"/home/nithin/workspace/go/service/",
	)
	if err != nil {
		slog.Error("Something went wrong while orchestrating", "err", err)
		return
	}

	// Step 2: Export to CSV
	csvt := CSVNodeExporter{}
	err = csvt.ExportTypes(extractRes.TypeDecls)
	if err != nil {
		slog.Error("", err)
	}

	err = csvt.ExportMembers(extractRes.Members)
	if err != nil {
		slog.Error("", err)
	}

	err = csvt.ExportInterfaces(extractRes.Interfaces)
	if err != nil {
		slog.Error("", err)
	}

	err = csvt.ExportNamed(extractRes.NamedTypes)
	if err != nil {
		slog.Error("", err)
	}

	err = csvt.ExportFile(extractRes.Files)
	if err != nil {
		slog.Error("", err)
	}

	err = csvt.ExportFunctions(extractRes.Functions)
	if err != nil {
		slog.Error("", err)
	}

	err = csvt.ExportNamespace(extractRes.Namespaces)
	if err != nil {
		slog.Error("", err)
	}

	csvr := CSVRelationshipExporter{}
	//
	err = csvr.ExportCalls(extractRes.Functions)
	if err != nil {
		slog.Error("", err)
	}

	err = csvr.ExportImplements(extractRes.TypeDecls)
	if err != nil {
		slog.Error("", err)
	}

	err = csvr.ExportImports(extractRes.Files)
	if err != nil {
		fmt.Println(err)
	}

	err = csvr.ExportReturns(extractRes.Functions)
	if err != nil {
		fmt.Println(err)
	}

}
