package extract

type ChunkType string

type Language string

const (
	Go   Language = "go"
	Java          = "java"
)

type Extractor interface {
	Extract(file string) ([]Chunk, error)
}

type Chunk struct {
	File        string
	Language    Language
	Content     string
	Type        ChunkType
	TokenCount  int
	ContentHash string
	BelongsTo   string
}

type GrepResult struct {
	Language string `json:"language"`
	Text     string `json:"text"`
	File     string `json:"file"`
}

const (
	Class        ChunkType = "class"
	Interface              = "interface"
	Method                 = "method"
	Function               = "function"
	LineComment            = "linecomment"
	BlockComment           = "blockcomment"
	Constant               = "constant"
	Struct                 = "struct"
	Import                 = "import"
	Package                = "package"
	Modifier               = "modifier"
	Constructor            = "constructor"
	Field                  = "field"
)
