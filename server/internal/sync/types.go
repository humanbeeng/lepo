package sync

type Syncer interface {
	Sync(target string) error
	Desync() error
}
type SyncRequest struct {
	URL string `json:"url" validate:"required"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}
