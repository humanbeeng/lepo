package process

import (
	"log/slog"

	"github.com/humanbeeng/lepo/prototypes/analysis/extract"
)

func Orchestrate(e extract.Extractor) {
	slog.Info("Begin orchestration")
	// Step 1: Extract
	extractRes, err := e.Extract(
		"github.com/humanbeeng/lepo/prototypes/go-testdata",
		"/Users/apple/workspace/go/lepo/prototypes/go-testdata",
	)
	if err != nil {
		slog.Error("", err)
	}

	// Step 2: Export to CSV
	csvt := CSVExporter{}

	err = csvt.ExportTypes(extractRes.TypeDecls)
	if err != nil {
		slog.Error("", err)
	}

	err = csvt.ExportInterfaces(extractRes.Interfaces)
	if err != nil {
		slog.Error("", err)
	}
}
