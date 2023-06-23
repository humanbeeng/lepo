package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/humanbeeng/lepo/server/internal/sync"
)

// TODO: Research more about project structure
type SyncControllerV1 struct {
	syncer sync.Syncer
}

var validate = validator.New()

func NewSyncControllerV1(syncer sync.Syncer) *SyncControllerV1 {
	return &SyncControllerV1{
		syncer: syncer,
	}
}

func AddV1SyncRoutes(router fiber.Router, controller *SyncControllerV1) {
	router.Post("/sync", controller.HandleSyncRequest)
	router.Post("/desync", controller.HandleDeSyncRequest)
}

func (s *SyncControllerV1) HandleSyncRequest(c *fiber.Ctx) error {
	req := new(sync.SyncRequest)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message)
	}

	errors := validateStruct(*req)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	err := s.syncer.Sync(req.URL)
	if err != nil {
		return err
	}
	return nil
}

func (s *SyncControllerV1) HandleDeSyncRequest(c *fiber.Ctx) error {
	req := new(sync.SyncRequest)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message)
	}

	errors := validateStruct(*req)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	err := s.syncer.Desync()
	if err != nil {
		return err
	}
	return nil
}

func validateStruct(req sync.SyncRequest) []*sync.ErrorResponse {
	var errors []*sync.ErrorResponse
	err := validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element sync.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
