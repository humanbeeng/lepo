package main

import (
	"fmt"
)

type Person struct {
	Name string
}

type Invoker interface {
	Invoke(int) Person
}

type Invokable struct{}

func (i Invokable) Invoke(n int) Person {
	return Person{Name: "Nithin"}
}

func DoSomething(i Invoker) {
	p := i.Invoke(1)
	fmt.Println(p.Name)
}

func main() {
	i := Invokable{}
	DoSomething(i)
}
