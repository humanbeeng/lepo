package golang

import (
	"bytes"
	"go/format"
	"go/token"
)

func extractCode(node any, fset *token.FileSet) (string, error) {
	var b []byte
	buf := bytes.NewBuffer(b)
	err := format.Node(buf, fset, node)
	if err != nil {
		return "", err
	}
	str := buf.String()

	return str, nil
}
