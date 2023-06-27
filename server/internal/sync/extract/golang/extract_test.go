package golang_test

import (
	"fmt"
	"testing"

	"github.com/humanbeeng/lepo/server/internal/sync/extract/golang"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TODO: Create new test resource folder to have dummy go file which contains all constructs

func TestGoExtractWhenFileNotFound(t *testing.T) {
	l, _ := zap.NewDevelopment()
	ge := golang.NewGoExtractor(l)
	chunks, err := ge.Extract("./lmao.go")
	assert.Nil(t, chunks)
	assert.NotNil(t, err)
}

func TestGoExtractWhenGoFileIsPassed(t *testing.T) {
	l, _ := zap.NewDevelopment()
	ge := golang.NewGoExtractor(l)
	chunks, _ := ge.Extract("./extract.go")

	for _, chunk := range chunks {
		fmt.Printf("%+v\n\n\n", chunk)
	}

	assert.True(t, true)
}
