package extract

type Language string

const (
	Go Language = "golang"
)

type Extractor interface {
	Extract(dirpath string) error
}

type GoExtractor struct {
	TypeDecls map[string]*TypeDecl
	Members   map[string]*Member
}

type Node struct {
	Name     string
	Code     string
	FilePath string
	Pos      int
	End      int
}

type Kind string

const (
	Interface Kind = "interface"
	Struct    Kind = "struct"
	Alias     Kind = "alias"
)

type Constant struct {
	Name       string
	QName      string
	TypeQName  string
	Underlying string
	Code       string
	Pos        int
	End        int
	Filepath   string
}

type TypeDecl struct {
	Name       string
	QName      string
	TypeQName  string
	Underlying string
	Code       string
	Kind       Kind
	Pos        int
	End        int
	Filepath   string
	// Package    string
}

type ExtractNodesResult struct{}

type File struct {
	Language Language
}

type Namespace struct {
	Node
}

type Member struct {
	Name        string
	QName       string
	TypeQName   string
	ParentQName string
	Code        string
	Pos         int
	End         int
	Filepath    string
}

type Function struct {
	Name         string
	QName        string
	ParentQName  string
	Code         string
	Pos          int
	End          int
	Filepath     string
	ReturnQNames []string
	ParamQNames  []string
}

func NewGoExtractor() *GoExtractor {
	return &GoExtractor{
		TypeDecls: make(map[string]*TypeDecl),
		Members:   make(map[string]*Member),
	}
}
