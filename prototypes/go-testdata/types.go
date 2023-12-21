package main

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/gofiber/fiber/v2"
	"github.com/humanbeeng/lepo/prototypes/go-testdata/testpackage"
)

type Visitor interface {
	Visit(ast.Node) ast.Visitor
}

type Person struct {
	Name      string
	Age       int
	TypesInfo types.Info
	Router    fiber.Route
}

func (p *Person) Display() {
	fmt.Println(p)
	somefunc()
	ft := testpackage.FuncTestStruct{}
	ft.Hello()
}

func somefunc() (string, error) {
	for i := 0; i < 10; i++ {
		message := anotherFunc()
		fmt.Println(message)
	}
	return "Hello there", nil
}

func anotherFunc() string {
	// There is nothing here
	return "From another func"
}

type MedicalData struct {
	Bloodtype string
}

type EducationData struct {
	Degree  string
	College string
}

type FinancialData struct {
	Networth int
	BankData struct {
		BankName string
		IFSC     int
		Anon     struct {
			innerAnonStruct struct {
				innerAnonName string
			}
		}
	}
}
