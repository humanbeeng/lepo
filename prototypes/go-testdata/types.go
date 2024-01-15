package main

import (
	"go/types"

	"github.com/gofiber/fiber/v2"
)

// Comment above struct decl
type Person struct {
	Name      string `json:"name"` // Inline comment
	Age       int    `json:"age"`
	TypesInfo types.Info
	Router    fiber.Route
}

// First line
// Second line
type MedicalData struct {
	// Field doc comment
	Bloodtype string // Field inline
}

/*
Block comment above struct
*/
type EducationData struct {
	Degree  string
	College string
}

type FinancialData struct {
	Networth int
	BankData struct {
		BankName string // Inner inline comment
		IFSC     int
	}
}
