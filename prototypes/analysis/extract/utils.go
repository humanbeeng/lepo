package extract

import (
	"bytes"
	"go/format"
	"go/token"
)

func code(node any, fset *token.FileSet) (string, error) {
	var codeStr string
	var b []byte
	buf := bytes.NewBuffer(b)
	err := format.Node(buf, fset, node)
	if err != nil {
		return "", err
	}

	codeStr = buf.String()
	return codeStr, nil
}
