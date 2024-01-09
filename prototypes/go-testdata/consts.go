package current

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type PowerRanger string // Type inline comment

// Type single line comment
type Success int

const DB string = "database_conn"

const (
	Red PowerRanger = "red"
	// First line
	// Second line
	Yellow PowerRanger = "yellow" // Third inline
	Blue   PowerRanger = "blue"
)

type Status int

const (
	Ok Status = iota // Const inline comment
	// Const single line comment
	Error
)

func (s *Status) Display() {
	fmt.Println(s)
	f := fiber.New()
	err := f.Listen("8080")
	if err != nil {
		panic(err)
	}
}
