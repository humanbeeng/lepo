package process

import (
	"strings"
)

func escapeStr(str string) string {
	// Replace double quotes with escaped double quotes
	str = strings.ReplaceAll(str, `"`, `\n`)
	// Replace newlines with escaped newlines
	str = strings.ReplaceAll(str, "\n", `\n`)
	// Replace backticks with escaped backticks
	str = strings.ReplaceAll(str, "`", "\\`")
	return str
}
