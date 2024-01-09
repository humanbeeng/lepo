package current

import (
	"log/slog"
)

type TestInterface interface {
	DoSomething(string) error
}

func AddStrings(a, b string) string {
	return a + b
}

func Invoke(msg string, p Printer) {
	printed, err := p.Print(msg)
	if err != nil {
		slog.Error("Unable to invoke printer", err)
	}
	slog.Info("Printed", "value", printed)
}
