package main

import (
	"go/types"

	"github.com/gofiber/fiber/v2"
)

type Person struct {
	Name      string
	Age       int
	TypesInfo types.Info
	Router    fiber.Route
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
	}
}
