package process

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"

	"github.com/humanbeeng/lepo/prototypes/analysis/extract"
)

type CSVRelationshipExporter struct{}

func (c *CSVRelationshipExporter) ExportImplements(types map[string]extract.TypeDecl) error {
	slog.Info("Exporting implements relationship to csv")

	csvFile, err := os.Create("neo4j/import/implements.csv")
	if err != nil {
		return fmt.Errorf("Unable to create implements.csv")
	}

	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	header := []string{
		"from",
		"to",
	}

	err = csvwriter.Write(header)
	if err != nil {
		return fmt.Errorf("Unable to write header to implements.csv")
	}

	for qname, t := range types {
		for _, impl := range t.ImplementsQName {
			row := []string{qname, impl}
			err := csvwriter.Write(row)
			if err != nil {
				return fmt.Errorf("unable to write row to implements.csv")
			}
		}
	}

	csvwriter.Flush()
	err = csvwriter.Error()
	if err != nil {
		return fmt.Errorf("unable to flush to implements.csv: %w", err)
	}

	slog.Info("Finished exporting to implements.csv")
	return nil
}

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
		return fmt.Errorf("unable to write header to calls.csv")
	}

	for qname, f := range functions {
		for _, call := range f.Calls {
			row := []string{qname, call}
			err := csvwriter.Write(row)
			if err != nil {
				return fmt.Errorf("unable to write call row to calls.csv: %w", err)
			}
		}
	}

	csvwriter.Flush()
	err = csvwriter.Error()
	if err != nil {
		return fmt.Errorf("unable to flush to calls.csv: %w", err)
	}
	slog.Info("Finished exporting callgraph to calls.csv")

	return nil
}

func (c *CSVRelationshipExporter) ExportImports(files map[string]extract.File) error {
	slog.Info("Exporting imports relationship to csv")

	csvFile, err := os.Create("neo4j/import/imports.csv")
	if err != nil {
		return fmt.Errorf("Unable to create imports.csv")
	}

	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	header := []string{
		"from",
		"to",
	}

	err = csvwriter.Write(header)
	if err != nil {
		return fmt.Errorf("unable to write header to imports.csv: %w", err)
	}

	for _, f := range files {
		for _, imports := range f.Imports {
			row := []string{f.Filename, imports.Path}
			err := csvwriter.Write(row)
			if err != nil {
				return fmt.Errorf("unable to write import row to imports.csv: %w", err)
			}
		}
	}

	csvwriter.Flush()
	err = csvwriter.Error()
	if err != nil {
		return fmt.Errorf("unable to flush to imports.csv: %w", err)
	}

	slog.Info("Finished exporting to imports.csv")
	return nil
}

func (c *CSVRelationshipExporter) ExportReturns(functions map[string]extract.Function) error {
	slog.Info("Exporting returns to csv")

	csvFile, err := os.Create("neo4j/import/returns.csv")
	if err != nil {
		return fmt.Errorf("Unable to create returns.csv")
	}

	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	header := []string{
		"from",
		"to",
	}

	err = csvwriter.Write(header)
	if err != nil {
		return fmt.Errorf("unable to write header to returns.csv")
	}

	for qname, f := range functions {
		for _, ret := range f.ReturnQNames {
			row := []string{qname, ret}
			err := csvwriter.Write(row)
			if err != nil {
				return fmt.Errorf("unable to write call row to returns.csv: %w", err)
			}
		}
	}

	csvwriter.Flush()
	err = csvwriter.Error()
	if err != nil {
		return fmt.Errorf("unable to flush to returns.csv: %w", err)
	}
	slog.Info("Finished exporting to returns.csv")

	return nil
}

func (c *CSVRelationshipExporter) ExportParams(functions map[string]extract.Function) error {
	slog.Info("Exporting params to csv")

	csvFile, err := os.Create("neo4j/import/params.csv")
	if err != nil {
		return fmt.Errorf("Unable to create params.csv")
	}

	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	header := []string{
		"from",
		"to",
	}

	err = csvwriter.Write(header)
	if err != nil {
		return fmt.Errorf("unable to write header to params.csv")
	}

	for qname, f := range functions {
		for _, param := range f.ParamQNames {
			row := []string{param, qname}
			err := csvwriter.Write(row)
			if err != nil {
				return fmt.Errorf("unable to write call row to param.csv: %w", err)
			}
		}
	}

	csvwriter.Flush()
	err = csvwriter.Error()
	if err != nil {
		return fmt.Errorf("unable to flush to param.csv: %w", err)
	}
	slog.Info("Finished exporting to param.csv")

	return nil
}
