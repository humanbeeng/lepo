package current

type PowerRanger string

const DB string = "database_conn"

const (
	Red    PowerRanger = "red"
	Yellow PowerRanger = "yellow"
	Blue   PowerRanger = "blue"
)

type Status int

const (
	Ok Status = iota
	Error
)
