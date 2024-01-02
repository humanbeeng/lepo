package current

import (
	"fmt"
)

type AnotherMethodStruct struct {
	Name string
}

func (a *AnotherMethodStruct) AnotherMethod() {
	fmt.Println(a.Name)
}
