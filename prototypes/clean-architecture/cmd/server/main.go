package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/humanbeeng/lepo/prototypes/clean-architecture/internal/command"
	"github.com/humanbeeng/lepo/prototypes/clean-architecture/internal/config"
	"github.com/humanbeeng/lepo/prototypes/clean-architecture/internal/database"
	storage "github.com/humanbeeng/lepo/prototypes/clean-architecture/internal/database"
)

func main() {
	var exitCode int

	defer func() {
		os.Exit(exitCode)
	}()

	appConfig, err := config.GetAppConfig()
	if err != nil {
		log.Printf("Unable to load app config file %v\n", err)
		exitCode = 1
		return
	}

	cleanup, err := run(*appConfig)
	defer cleanup()
	if err != nil {
		log.Printf("Error while starting the server %v\n", err)
		exitCode = 1
		return
	}

	command.Gracefully()
}

func run(appConfig config.AppConfig) (func(), error) {
	app, cleanup, err := buildServer(appConfig)
	if err != nil {
		return nil, err
	}

	go func() {
		app.Listen(":" + fmt.Sprintf("%d", appConfig.ServerConfig.Port))
	}()

	return func() {
		cleanup()
		app.Shutdown()
	}, nil
}

func buildServer(appConfig config.AppConfig) (*fiber.App, func(), error) {
	db, err := storage.BootstrapMySQL()
	if err != nil {
		return nil, nil, err
	}
	app := fiber.New()
	app.Use(logger.New())

	app.Get("/internal/health", func(c *fiber.Ctx) error {
		return c.SendString("I'm healthy !")
	})

	return app, func() {
		database.CloseMySQLConnection(db)
	}, nil
}
