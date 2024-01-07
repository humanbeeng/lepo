package current

import (
	"fmt"
)

type AnotherMethodStruct struct {
	Name string
}

// This is a comment
func (a *AnotherMethodStruct) AnotherMethod() {
	// This is a body comment
	fmt.Println(a.Name)
}
