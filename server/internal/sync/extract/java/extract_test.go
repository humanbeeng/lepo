package java_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/humanbeeng/lepo/server/internal/sync/extract/java"
)

func TestHello(t *testing.T) {
	assert.True(t, true)
}

func TestGoTreesitter(t *testing.T) {
	l, _ := zap.NewDevelopment()
	e := java.NewJavaExtractor(l)

	_, err := e.Extract(
		"/home/humanbeeng/projects/lepo/prototypes/extract/java/src/main/java/ai/lepo/java/ImplementationClass.java",
	)
	assert.Nil(t, err)
}
