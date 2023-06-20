package sync

type Syncer interface {
	Sync(path string) error
}
