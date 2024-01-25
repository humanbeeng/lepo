package main

import (
	"github.com/humanbeeng/lepo/prototypes/analysis/extract/golang"
	"github.com/humanbeeng/lepo/prototypes/analysis/process"
)

func main() {
	// TODO : Remove extract interface and pass cmd line args for dir and pkg
	e := golang.NewGoExtractor()
	process.Orchestrate(e)
}
