package main

import "github.com/humanbeeng/lepo/prototypes/extract/internal/sync"

func main() {
	syncOpts := sync.DirectorySyncerOpts{
		ExcludedFolderPatterns: make([]string, 0),
	}

	syncer := sync.NewDirectorySyncer(syncOpts)

	syncer.Sync(".")
}
