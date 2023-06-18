package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/humanbeeng/lepo/server/internal/sync"
)

// TODO: Research more about project structure
type SyncRouter struct{}

type SyncRequestBody struct {
	URL string
}

func NewSyncRouter(router fiber.Router) {
	router.Post("/sync", RequestSyncHandler)
}

func RequestSyncHandler(c *fiber.Ctx) error {
	var syncRequestBody SyncRequestBody

	if err := c.BodyParser(&syncRequestBody); err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message)
	}

	opts := sync.GitSyncerOpts{
		URL: syncRequestBody.URL,
	}
	syncer := sync.NewGitSyncer(opts)
	err := syncer.Sync()
	if err != nil {
		return err
	}
	return nil
}
