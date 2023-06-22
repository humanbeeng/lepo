package main

import (
	_ "github.com/humanbeeng/lepo/prototypes/extract/internal/embed"
	"github.com/humanbeeng/lepo/prototypes/extract/internal/sync"
)

func main() {
	syncOpts := sync.DirectorySyncerOpts{
		ExcludedFolderPatterns: make([]string, 0),
	}

	syncer := sync.NewDirectorySyncer(syncOpts)
	folder := "/home/personal/projects/go/go-lb/"

	syncer.Sync(folder)
}
