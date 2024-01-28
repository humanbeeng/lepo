package current

import (
	"fmt"
)

type Printer func(msg string) string

type NiceGreeter struct {
	Name string
}

func (nc *NiceGreeter) PrintHello(msg string) error {
	fmt.Println("Hello")
	return nil
}

func (nc *NiceGreeter) Greet(msg string) (string, error) {
	fmt.Println("Greeting nicely", msg)
	var p Printer
	greetMsg := p("hello there")
	fmt.Println(greetMsg)
	return msg, nil
}
