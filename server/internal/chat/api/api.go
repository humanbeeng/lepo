package api

import "github.com/gofiber/fiber/v2"

type ChatRouter struct{}

func NewChatRouter(app fiber.Router) {
	app.Get("/chat", GetChatsHandler)
	app.Post("/chat", PostChatHandler)
}

func GetChatsHandler(c *fiber.Ctx) error {
	return c.SendString("Hello from chat endpoint")
}

func PostChatHandler(c *fiber.Ctx) error {
	return c.SendString("Posting something")
}
