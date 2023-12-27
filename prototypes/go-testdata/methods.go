package main

type MethodTestStruct struct {
	StringMember string
}

type Opts struct {
	num int
}

func (mts *MethodTestStruct) MethodOne(a Opts) error {
	// fmt.Println(a.num)
	// res, _ := add(1, 2)
	//
	// fmt.Println("Result", res)

	// s1 := "hello"
	// s2 := " there"

	// concatStr := internal.AddStrings(s1, s2)
	// fmt.Println(concatStr)
	recursive(10)

	return nil
}

// func add(a, b int) (int, error) {
// 	return (a + b), nil
// }

func recursive(a int) int {
	if a <= 0 {
		return a
	}
	return recursive(a - 1)
}

//
// func handleError(err error) {
// 	fmt.Println(err)
// }
//
// func returnOpts() Opts {
// 	return Opts{num: 1}
// }
//
// func returnNothing() {
// 	go func(a int) {
// 		fmt.Println(a)
// 	}(2)
// 	fmt.Println("I return nothing")
// }
