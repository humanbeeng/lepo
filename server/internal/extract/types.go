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
	Module      string
}

type GrepResult struct {
	Language string `json:"language"`
	Text     string `json:"text"`
	File     string `json:"file"`
}

const (
	Class        ChunkType = "class"
	Interface    ChunkType = "interface"
	Method       ChunkType = "method"
	Function     ChunkType = "function"
	LineComment  ChunkType = "linecomment"
	BlockComment ChunkType = "blockcomment"
	Constant     ChunkType = "constant"
	Struct       ChunkType = "struct"
	Import       ChunkType = "import"
	Package      ChunkType = "package"
	Modifier     ChunkType = "modifier"
	Constructor  ChunkType = "constructor"
	Field        ChunkType = "field"
)
