package process

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"

	"github.com/humanbeeng/lepo/prototypes/analysis/extract"
)

type (
	CSVNodeExporter struct{}
)

func (c *CSVNodeExporter) ExportTypes(types map[string]extract.TypeDecl) error {
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

func (c *CSVNodeExporter) ExportFunctions(functions map[string]extract.Function) error {
	slog.Info("Exporting functions to csv")

	csvFile, err := os.Create("neo4j/import/function.csv")
	if err != nil {
		return fmt.Errorf("Unable to create function.csv %v", err)
	}

	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	header := []string{
		extract.Name,
		extract.QualifiedName,
		extract.ParentQualifiedName,
		"kind",
		extract.Code,
		"doc",
		"file",
	}

	err = csvwriter.Write(header)
	if err != nil {
		return fmt.Errorf("Unable to write header to function.csv %v", err)
	}

	for qname, f := range functions {
		code := f.Code

		code = escapeStr(code)

		row := []string{f.Name, qname, f.ParentQName, "function", code, f.Doc.Comment, f.Filepath}
		err := csvwriter.Write(row)
		if err != nil {
			slog.Error("Unable to write to function.csv file", err)
			return err
		}
	}

	csvwriter.Flush()
	err = csvwriter.Error()
	if err != nil {
		slog.Error("Unable to flush to function.csv file", err)
	}
	slog.Info("Finished exporting functions to function.csv file")

	return nil
}

func (c *CSVNodeExporter) ExportInterfaces(types map[string]extract.TypeDecl) error {
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

func (c *CSVNodeExporter) ExportNamed(named map[string]extract.Named) error {
	slog.Info("Exporting named types to csv")

	csvFile, err := os.Create("neo4j/import/named.csv")
	if err != nil {
		return fmt.Errorf("Unable to create named.csv %v", err)
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
		return fmt.Errorf("Unable to write header to named.csv %v", err)
	}

	for qname, n := range named {
		tqn := n.TypeQName
		und := n.Underlying
		code := n.Code

		code = escapeStr(code)
		und = escapeStr(und)
		tqn = escapeStr(tqn)

		// TODO : Refer schema for named types and revisit this name: named
		row := []string{n.Name, qname, tqn, und, "named", code, n.Doc.Comment}
		err := csvwriter.Write(row)
		if err != nil {
			slog.Error("Unable to write to named.csv file", err)
			return err
		}
	}

	csvwriter.Flush()
	err = csvwriter.Error()
	if err != nil {
		slog.Error("Unable to flush to named.csv file", err)
	}
	slog.Info("Finished exporting types to named.csv file")

	return nil
}

func (c *CSVNodeExporter) ExportFile(files map[string]extract.File) error {
	slog.Info("Exporting go file information as csv")

	filecsv, err := os.Create("neo4j/import/file.csv")
	if err != nil {
		return fmt.Errorf("Unable to create file.csv %v", err)
	}

	importcsv, err := os.Create("neo4j/import/import.csv")
	if err != nil {
		return fmt.Errorf("Unable to create import.csv %v", err)
	}

	defer filecsv.Close()
	defer importcsv.Close()

	fileWriter := csv.NewWriter(filecsv)
	importWriter := csv.NewWriter(importcsv)

	fileHeader := []string{
		extract.Filename,
		extract.Package,
		extract.Language,
	}

	importHeader := []string{
		extract.Name,
		extract.Path,
		extract.Comment,
	}

	err = fileWriter.Write(fileHeader)
	if err != nil {
		// TODO: Refactor slog error. Follow best practice
		return fmt.Errorf("Unable to write header to file.csv %v", err)
	}

	err = importWriter.Write(importHeader)
	if err != nil {
		return fmt.Errorf("Unable to write header to import.csv %v", err)
	}

	for _, file := range files {
		// TODO : Refer schema for named types and revisit this name: named
		row := []string{file.Filename, file.Package, file.Language}
		err := fileWriter.Write(row)
		if err != nil {
			slog.Error("Unable to write row to file.csv file", err)
			return err
		}

		for _, i := range file.Imports {
			importRow := []string{i.Name, i.Path, i.Doc.Comment}
			err := importWriter.Write(importRow)
			if err != nil {
				slog.Error("Unable to write row to import.csv", err)
				return err
			}
		}
	}

	fileWriter.Flush()
	err = fileWriter.Error()
	if err != nil {
		slog.Error("Unable to flush to file.csv file", err)
	}

	importWriter.Flush()
	err = importWriter.Error()
	if err != nil {
		slog.Error("Unable to flush to import.csv file", err)
	}

	slog.Info("Finished exporting files to file.csv and import.csv")

	return nil
}

func (c *CSVNodeExporter) ExportNamespace(namespaces []extract.Namespace) error {
	slog.Info("Exporting structs to csv")

	csvFile, err := os.Create("neo4j/import/namespace.csv")
	if err != nil {
		return fmt.Errorf("Unable to create namespace.csv %v", err)
	}

	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	header := []string{
		extract.Name,
	}

	err = csvwriter.Write(header)
	if err != nil {
		return fmt.Errorf("Unable to write header to namespace.csv %v", err)
	}

	for _, t := range namespaces {

		row := []string{t.Name}
		err := csvwriter.Write(row)
		if err != nil {
			slog.Error("Unable to write to namespace.csv file", err)
			return err
		}
	}

	csvwriter.Flush()
	err = csvwriter.Error()
	if err != nil {
		slog.Error("Unable to flush to namespace.csv file", err)
	}

	slog.Info("Finished exporting types to namespace.csv file")
	return nil
}

func (c *CSVNodeExporter) ExportMembers(members map[string]extract.Member) error {
	slog.Info("Exporting members to csv")

	csvFile, err := os.Create("neo4j/import/members.csv")
	if err != nil {
		return fmt.Errorf("Unable to create members.csv %v", err)
	}

	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	header := []string{
		extract.Name,
		extract.QualifiedName,
		extract.TypeName,
		extract.ParentQualifiedName,
		"kind",
		extract.Code,
		"doc",
	}

	err = csvwriter.Write(header)
	if err != nil {
		// TODO: Refactor slog error. Follow best practice
		return fmt.Errorf("Unable to write header to members.csv %v", err)
	}

	for qname, t := range members {
		tqn := t.TypeQName
		code := t.Code

		code = escapeStr(code)
		tqn = escapeStr(tqn)

		row := []string{
			t.Name,
			qname,
			tqn,
			t.ParentQName,
			// TODO: Refactor kind as header constant
			"kind",
			code,
			t.Doc.Comment,
		}
		err := csvwriter.Write(row)
		if err != nil {
			slog.Error("Unable to write to members.csv file", err)
			return err
		}
	}

	csvwriter.Flush()
	err = csvwriter.Error()
	if err != nil {
		slog.Error("Unable to flush to members.csv file", err)
	}

	slog.Info("Finished exporting types to members.csv file")
	return nil
}
