package extract

type Language string

const (
	Go Language = "golang"
)

type Extractor interface {
	Extract(dirpath string) error
}

type GoExtractor struct {
	TypeDefs map[string]*TypeDef
	Members  map[string]*Member
}

type Node struct {
	Name          string
	QualifiedName string
	Code          string
	Pos           int
	End           int
	File          string
	Package       string
	Visibility    Visibility
}

type Kind string

type Visibility string

const (
	Public  Visibility = "public"
	Private Visibility = "private"
	Package Visibility = "package"
)

const (
	Interface Kind = "interface"
	Struct    Kind = "struct"
	Alias     Kind = "alias"
)

type TypeDef struct {
	Node
	Type       string
	Underlying string
	Kind       Kind
	Alias      string
	AliasFor   string
	Satisfies  map[string]byte // all keys are always be qualified name
	DependsOn  map[string]byte
	Signatures map[string]byte
	Types      map[string]byte
	Methods    map[string]byte
	Members    map[string]byte
}

type File struct {
	Language Language
}

type Member struct {
	Node
	Type   string
	Parent string
}

type Namespace struct {
	Node
}

type Function struct {
	Node
	Parent    string
	Overrides string
	Returns   map[string]byte
	Params    map[string]byte
}

func NewGoExtractor() *GoExtractor {
	return &GoExtractor{
		TypeDefs: make(map[string]*TypeDef),
		Members:  make(map[string]*Member),
	}
}
