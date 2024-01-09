package current

import (
	"fmt"
	"log/slog"

	f "github.com/gofiber/fiber/v2" // Import comment
)

type MethodTestStruct struct {
	StringMember string
}

type Opts struct {
	Num int
}

func (mts *MethodTestStruct) MethodOne(a Opts) error {
	returnNothing()
	res := AddStrings("Hello", "There")
	f.New()
	slog.Info("Added two strings", "result", res)
	cp := ConsolePrinter{}
	Invoke("Hello there", cp)
	return nil
}

func returnNothing() {
	go func(a int) {
		fmt.Println(a)
	}(2)
	fmt.Println("I return nothing")
}
