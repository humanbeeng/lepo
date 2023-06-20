package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	chatapi "github.com/lepoai/lepo/server/internal/chat/api"
	"github.com/lepoai/lepo/server/internal/command"
	config "github.com/lepoai/lepo/server/internal/config"
	"github.com/lepoai/lepo/server/internal/database"
	storage "github.com/lepoai/lepo/server/internal/database"
	syncapi "github.com/lepoai/lepo/server/internal/sync/api"
	tempapi "github.com/lepoai/lepo/server/internal/temp/api"
)

func main() {
	var exitCode int

	defer func() {
		os.Exit(exitCode)
	}()

	err := config.LoadAppConfig()
	if err != nil {
		log.Printf("error: %v\n", err)
		exitCode = 1
		return
	}

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
	app, cleanup, err := initComponents(appConfig)
	if err != nil {
		return nil, err
	}

	go func() {
		err = app.Listen(":" + fmt.Sprintf("%d", appConfig.ServerConfig.Port))
		if err != nil {
			panic(err)
		}
	}()

	return func() {
		cleanup()
		err = app.Shutdown()
		if err != nil {
			panic(err)
		}
	}, nil
}

func initComponents(appConfig config.AppConfig) (*fiber.App, func(), error) {
	db, err := storage.BootstrapMySQL()
	if err != nil {
		return nil, nil, err
	}

	_, err = storage.BootStrapWeaviate()
	if err != nil {
		return nil, nil, err
	}

	app := fiber.New()
	addRoutes(app)

	return app, func() {
		err = database.CloseMySQLConnection(db)
		if err != nil {
			panic(err)
		}
	}, nil
}

func addRoutes(app *fiber.App) {
	app.Use(logger.New())
	v1 := app.Group("/v1")
	v1.Get("/_internal/health", func(c *fiber.Ctx) error {
		return c.SendString("I'm healthy !")
	})
	chatapi.NewChatRouter(v1)
	syncapi.NewSyncRouter(v1)
	tempapi.NewTempRouter(v1)
}
