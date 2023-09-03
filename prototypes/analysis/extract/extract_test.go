package extract

import (
	"testing"
)

func TestGoExtract(t *testing.T) {
	e := NewGoExtractor()
	e.Extract("github.com/humanbeeng/lepo/prototypes/go-testdata")
}
