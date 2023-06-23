package v1

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/humanbeeng/lepo/server/internal/chat"
)

type ChatControllerV1 struct {
	resolver chat.ChatResolver
}

func NewChatControllerV1(resolver chat.ChatResolver) *ChatControllerV1 {
	return &ChatControllerV1{
		resolver: resolver,
	}
}

func AddV1ChatRoutes(router fiber.Router, controller *ChatControllerV1) {
	router.Post("/chat", controller.PostChat)
}

func (cont *ChatControllerV1) PostChat(c *fiber.Ctx) error {
	var req chat.ChatRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message)
	}

	errors := validateStruct(req)

	if errors != nil {
		for _, e := range errors {
			fmt.Println(e.Tag)
		}
		return fiber.NewError(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message)
	}

	fmt.Printf("Received request %+v\n", req)

	resp, err := cont.resolver.Resolve(req)
	if err != nil {
		log.Println("Unable to resolve", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	r, _ := json.Marshal(resp)

	return c.Send(r)
}

var validate = validator.New()

func validateStruct(req chat.ChatRequest) []*chat.ErrorResponse {
	var errors []*chat.ErrorResponse
	err := validate.Struct(req)
	// Revisit: Check if this elaboration is needed ?
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element chat.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
