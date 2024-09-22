package extract

// TODO: rename package into namespace

// Languages
const (
	Go         string = "golang"
	Rust       string = "rust"
	Java       string = "java"
	JavaScript string = "javascript"
	TypeScript string = "typescript"
)

type ExtractNodesResult struct {
	TypeDecls  map[string]TypeDecl
	Members    map[string]Member
	Interfaces map[string]TypeDecl
	Functions  map[string]Function
	NamedTypes map[string]Named
	Files      map[string]File
	Namespaces []Namespace
	Vars       map[string]Variable
}

type Extractor interface {
	Extract(pkgstr string, dir string) (ExtractNodesResult, error)
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

const (
	Name                string = "name"
	QualifiedName       string = "qualified_name"
	TypeName            string = "type"
	UnderlyingType      string = "underlying_type"
	ParentQualifiedName string = "parent_qualified_name"
	Implements          string = "implements"
	Code                string = "code"
	Filename            string = "file"
	Package             string = "package"
	Language            string = "language"
	Path                string = "path"
	Comment             string = "comment"
)

type Variable struct {
	Name       string
	QName      string
	TypeQName  string
	Underlying string

	Code      string
	Doc       Doc
	Pos       int
	End       int
	Filepath  string
	Namespace string
}

// TODO: Refactor this as TypeDecl
type Named struct {
	Name       string
	QName      string
	Namespace  Namespace
	TypeQName  string
	Underlying string
	Code       string
	Doc        Doc
	Pos        int
	End        int
	Filepath   string
}

type TypeDecl struct {
	Name            string
	QName           string
	TypeQName       string
	Underlying      string
	ImplementsQName []string
	Code            string
	Doc             Doc
	Kind            Kind
	Pos             int
	End             int
	Filepath        string
	Namespace       Namespace
}

type File struct {
	Filename string
	// package this file belongs to.
	Namespace string
	Imports   []Import
	Language  string
}

type Namespace struct {
	Name string
}

type Member struct {
	Name        string
	QName       string
	Namespace   Namespace
	TypeQName   string
	ParentQName string
	Code        string
	Doc         Doc
	Pos         int
	End         int
	Filepath    string
}

type Function struct {
	Name         string
	QName        string
	Namespace    Namespace
	ParentQName  string
	Calls        []string
	Doc          Doc
	Code         string
	Pos          int
	End          int
	Filepath     string
	ReturnQNames []string
	ParamQNames  []string
}

type DocType byte

const (
	SingleLine DocType = iota
	MultiLine
	Block
	Inline
)

type Import struct {
	Name string
	Path string
	Doc  Doc
}

type Doc struct {
	Comment string
	OfQName string
	Type    DocType
}
