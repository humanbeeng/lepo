package sync_test

import (
	"testing"

	"github.com/humanbeeng/lepo/prototypes/extract/internal/sync"
	"github.com/stretchr/testify/assert"
)

func TestDirectorySyncForValidGoTarget(t *testing.T) {
	opts := sync.DirectorySyncerOpts{
		ExcludedFolderPatterns: make([]string, 0),
	}

	dirSyncer := sync.NewDirectorySyncer(opts)

	err := dirSyncer.Sync(".")

	assert.Nil(t, err)
}

func TestDirectorySyncForInValidTarget(t *testing.T) {
	opts := sync.DirectorySyncerOpts{
		ExcludedFolderPatterns: make([]string, 0),
	}

	dirSyncer := sync.NewDirectorySyncer(opts)

	err := dirSyncer.Sync("non-existing-folder")
	assert.NotNil(t, err)
}
