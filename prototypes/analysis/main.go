package main

import (
	"github.com/humanbeeng/lepo/prototypes/analysis/extract/golang"
	"github.com/humanbeeng/lepo/prototypes/analysis/process"
)

func main() {
	e := golang.NewGoExtractor()
	process.Orchestrate(e)
}
