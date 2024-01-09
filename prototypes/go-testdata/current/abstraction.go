package current

import (
	"fmt"
)

type Printer interface {
	Print(string) (string, error)
}

type ConsolePrinter struct{}

func (c ConsolePrinter) Print(msg string) (string, error) {
	fmt.Println("Writing to console")
	return msg, nil
}
