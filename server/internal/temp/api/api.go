package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/lepoai/lepo/server/internal/database"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

type TempRouter struct{}

func NewTempRouter(router fiber.Router) {
	router.Post("/temp/search", SearchHandler)
}

type SearchRequest struct {
	Query string
}

func SearchHandler(c *fiber.Ctx) error {
	var r SearchRequest

	if err := c.BodyParser(&r); err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message)
	}

	fields := []graphql.Field{
		{Name: "code"},
		{Name: "language"},
		{Name: "module"},
		{Name: "file"},
		{Name: "codeType"},
	}

	// TODO: Remove this query and move this into a test maybe ?
	nearText := database.WeaviateClient.GraphQL().
		NearTextArgBuilder().
		WithConcepts([]string{r.Query})

	result, err := database.WeaviateClient.GraphQL().Get().
		WithClassName("CodeSnippets").
		WithLimit(10).
		WithNearText(nearText).
		WithFields(fields...).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	data := result.Data

	fmt.Println(data)

	m, _ := json.Marshal(data)

	return c.Send(m)
}
