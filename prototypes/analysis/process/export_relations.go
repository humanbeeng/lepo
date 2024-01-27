package process

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"

	"github.com/humanbeeng/lepo/prototypes/analysis/extract"
)

type CSVRelationshipExporter struct{}

func (c *CSVRelationshipExporter) ExportCalls(functions map[string]extract.Function) error {
	slog.Info("Exporting calls to csv")

	csvFile, err := os.Create("neo4j/import/calls.csv")
	if err != nil {
		return fmt.Errorf("Unable to create calls.csv")
	}

	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	header := []string{
		"from",
		"to",
	}

	err = csvwriter.Write(header)
	if err != nil {
		return fmt.Errorf("Unable to write header to calls.csv")
	}

	for qname, f := range functions {
		for _, call := range f.Calls {
			row := []string{qname, call}
			err := csvwriter.Write(row)
			if err != nil {
				slog.Error("Unable to write call row to calls.csv", err)
				return err
			}
		}
	}

	csvwriter.Flush()
	err = csvwriter.Error()
	if err != nil {
		slog.Error("Unable to flush to calls.csv", err)
	}
	slog.Info("Finished exporting callgraph to calls.csv")

	return nil
}
