package main

import "github.com/humanbeeng/lepo/prototypes/analysis/extract"

func main() {
	e := extract.NewGoExtractor()
	invoke(e)
}

func invoke(e extract.Extractor) {
	err := e.Extract("github.com/dgraph-io/dgraph")
	if err != nil {
		panic(err)
	}
}
