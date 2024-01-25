package process

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"

	"github.com/humanbeeng/lepo/prototypes/analysis/extract"
)

type CSVExporter struct{}

func (c *CSVExporter) ExportTypes(types map[string]extract.TypeDecl) error {
	slog.Info("Exporting structs to csv")

	csvFile, err := os.Create("neo4j/import/type.csv")
	if err != nil {
		return fmt.Errorf("Unable to create type.csv %v", err)
	}

	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	header := []string{
		extract.Name,
		extract.QualifiedName,
		extract.TypeName,
		extract.UnderlyingType,
		"kind",
		extract.Code,
		"doc",
	}

	err = csvwriter.Write(header)
	if err != nil {
		// TODO: Refactor slog error. Follow best practice
		return fmt.Errorf("Unable to write header to type.csv %v", err)
	}

	for qname, t := range types {
		tqn := t.TypeQName
		und := t.Underlying
		code := t.Code

		code = escapeStr(code)
		und = escapeStr(und)
		tqn = escapeStr(tqn)

		row := []string{t.Name, qname, tqn, und, string(t.Kind), code, t.Doc.Comment}
		err := csvwriter.Write(row)
		if err != nil {
			slog.Error("Unable to write to type.csv file", err)
			return err
		}
	}
	csvwriter.Flush()
	err = csvwriter.Error()
	if err != nil {
		slog.Error("Unable to flush to type.csv file", err)
	}

	slog.Info("Finished exporting types to type.csv file")
	return nil
}

func (c *CSVExporter) ExportInterfaces(types map[string]extract.TypeDecl) error {
	slog.Info("Exporting interfaces to csv")

	csvFile, err := os.Create("neo4j/import/interface.csv")
	if err != nil {
		return fmt.Errorf("Unable to create interface.csv %v", err)
	}

	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	header := []string{
		extract.Name,
		extract.QualifiedName,
		extract.TypeName,
		extract.UnderlyingType,
		"kind",
		extract.Code,
		"doc",
	}

	err = csvwriter.Write(header)
	if err != nil {
		// TODO: Refactor slog error. Follow best practice
		return fmt.Errorf("Unable to write header to interface.csv %v", err)
	}

	for qname, t := range types {
		tqn := t.TypeQName
		und := t.Underlying
		code := t.Code

		code = escapeStr(code)
		und = escapeStr(und)
		tqn = escapeStr(tqn)

		row := []string{t.Name, qname, tqn, und, string(t.Kind), code, t.Doc.Comment}
		err := csvwriter.Write(row)
		if err != nil {
			slog.Error("Unable to write to interface.csv file", err)
			return err
		}
	}

	csvwriter.Flush()
	err = csvwriter.Error()
	if err != nil {
		slog.Error("Unable to flush to interface.csv file", err)
	}
	slog.Info("Finished exporting types to interface.csv file")

	return nil
}
