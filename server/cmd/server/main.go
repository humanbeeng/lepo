package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	chatapi "github.com/humanbeeng/lepo/server/internal/chat/api/v1"
	"github.com/humanbeeng/lepo/server/internal/chat/function"
	"github.com/humanbeeng/lepo/server/internal/chat/resolver"
	config "github.com/humanbeeng/lepo/server/internal/config"
	"github.com/humanbeeng/lepo/server/internal/database"
	storage "github.com/humanbeeng/lepo/server/internal/database"
	"github.com/humanbeeng/lepo/server/internal/git"
	"github.com/humanbeeng/lepo/server/internal/sync"
	syncapi "github.com/humanbeeng/lepo/server/internal/sync/api/v1"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
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

	cleanup, err := runApp(*appConfig)
	defer cleanup()
	if err != nil {
		log.Printf("Error while starting the server %v\n", err)
		exitCode = 1
		return
	}

	gracefully()
}

func runApp(appConfig config.AppConfig) (func(), error) {
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
	zap, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	db, err := storage.BootstrapMySQL()
	if err != nil {
		return nil, nil, err
	}

	wvt, err := storage.BootStrapWeaviate()
	if err != nil {
		return nil, nil, err
	}

	openai := openai.NewClient(appConfig.OpenAIConfig.ApiKey)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	v1 := app.Group("/v1")
	v1.Get("/_internal/health", func(c *fiber.Ctx) error {
		return c.SendString("I'm healthy !")
	})

	codeCtxFunc := function.NewCodeContextFunction(wvt)

	resolver := resolver.NewOpenAIChatResolver(openai, codeCtxFunc, zap)

	chatController := chatapi.NewChatControllerV1(resolver)
	chatapi.AddV1ChatRoutes(v1, chatController)

	cloner := git.NewGitCloner()
	syncer := sync.NewGitSyncer(cloner, wvt)
	syncController := syncapi.NewSyncControllerV1(syncer)
	syncapi.AddV1SyncRoutes(v1, syncController)

	return app, func() {
		err = database.CloseMySQLConnection(db)
		if err != nil {
			panic(err)
		}
	}, nil
}

func gracefully() {
	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("App Shutdown Gracefully")
}
