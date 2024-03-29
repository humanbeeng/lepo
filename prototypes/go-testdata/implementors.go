package main

import (
	"fmt"

	"github.com/humanbeeng/lepo/prototypes/go-testdata/current"
)

type Greeter interface {
	Greet(string) (string, error)
}

type HelloPrinter interface {
	PrintHello(string) error
}

type FancyGreeter struct{}

func (fc FancyGreeter) Greet(msg string) (string, error) {
	fmt.Println(msg)
	return msg, fmt.Errorf("Nothing")
}

type RudeGreeter struct{}

func (rg RudeGreeter) Greet(msg string) (string, error) {
	return msg, nil
}

func DoSomething(g Greeter) {
	_, err := g.Greet("GG")
	if err != nil {
		fmt.Println(err)
	}
}

func DoAnother(hp HelloPrinter) {
	hp.PrintHello("hello")
}

func main() {
	// fc := FancyGreeter{}
	rg := RudeGreeter{}
	ng := &current.NiceGreeter{}
	DoAnother(ng)
	DoSomething(rg)
	DoSomething(ng)
	// DoSomething(fc)
}
