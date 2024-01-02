package current

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type PowerRanger string

const DB string = "database_conn"

const (
	Red    PowerRanger = "red"
	Yellow PowerRanger = "yellow"
	Blue   PowerRanger = "blue"
)

type Status int

const (
	Ok Status = iota
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
