package main

import "fmt"

type MethodTestStruct struct {
	StringMember string
}

type Opts struct {
	num int
}

func (mts *MethodTestStruct) MethodOne(a Opts) error {
	fmt.Println(a.num)
	res, err := add(1, 2)
	if err != nil {
		handleError(err)
	}
	fmt.Println("Result", res)
	return nil
}

func add(a, b int) (int, error) {
	return (a + b), nil
}

func handleError(err error) {
	fmt.Println(err)
}

func returnOpts() Opts {
	return Opts{num: 1}
}

func returnNothing() {
	go func(a int) {
		fmt.Println(a)
	}(2)
	fmt.Println("I return nothing")
}
