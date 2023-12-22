package main

import "fmt"

type MethodTestStruct struct {
	StringMember string
}

func (mts *MethodTestStruct) MethodOne() {
	a, err := add(1, 2)
	if err != nil {
		handleError(err)
	}
	fmt.Println("Result", a)
}

func add(a, b int) (int, error) {
	return (a + b), nil
}

func handleError(err error) {
	fmt.Println(err)
}
