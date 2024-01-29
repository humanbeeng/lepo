package process

import (
	"log/slog"

	"github.com/humanbeeng/lepo/prototypes/analysis/extract"
)

func Orchestrate(e extract.Extractor) {
	slog.Info("Begin orchestration")
	// Step 1: Extract
	extractRes, err := e.Extract(
		// "github.com/humanbeeng/lepo/prototypes/go-testdata",
		// "/Users/apple/workspace/go/lepo/prototypes/go-testdata",

		"github.com/dgraph-io/dgraph",
		"/Users/apple/workspace/misc/dgraph",
	)
	if err != nil {
		slog.Error("", err)
	}

	// Step 2: Export to CSV
	csvt := CSVNodeExporter{}
	err = csvt.ExportTypes(extractRes.TypeDecls)
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

	err = csvr.ExportCalls(extractRes.Functions)
	if err != nil {
		slog.Error("", err)
	}
}
