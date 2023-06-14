package chunk

type Chunk struct {
	Filename  string
	Content   string
	Language  Language
	DocString string
	Type
}

type FileChunker interface {
	chunk(file string) ([]Chunk, error)
}

type Language uint

type Type uint

const (
	Method Type = iota
	Struct
	Interface
)

const (
	Go Language = iota
)
