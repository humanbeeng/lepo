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
	csvFile, err := os.Create("neo4j/import/go-types.csv")
	if err != nil {
		return fmt.Errorf("Unable to create go-types.csv %v", err)
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	header := []string{"name", "qname", "typeqname", "underlying", "code", "doc"}

	err = csvwriter.Write(header)
	if err != nil {
		// TODO: Refactor slog error. Follow best practice
		return fmt.Errorf("Unable to write header to go-types.csv %v", err)
	}

	for k, v := range types {
		tqn := v.TypeQName
		und := v.Underlying
		code := v.Code

		code = escapeStr(code)
		und = escapeStr(und)
		tqn = escapeStr(tqn)

		row := []string{v.Name, k, tqn, und, code, v.Doc.Comment}
		err := csvwriter.Write(row)
		if err != nil {
			slog.Error("Unable to write to CSV file", err)
			return err
		}
	}

	csvwriter.Flush()
	return nil
}
