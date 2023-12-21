package main

import (
	"github.com/humanbeeng/lepo/prototypes/analysis/extract"
)

func main() {
	e := extract.NewGoExtractor()
	err := e.Extract("github.com/humanbeeng/lepo/prototypes/analysis")
	if err != nil {
		panic(err)
	}
}
