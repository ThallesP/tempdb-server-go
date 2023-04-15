package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bmdavis419/tapir-app/config"
	_ "github.com/bmdavis419/tapir-app/docs"
	"github.com/bmdavis419/tapir-app/internal/databases"
	"github.com/bmdavis419/tapir-app/internal/storage"
	"github.com/bmdavis419/tapir-app/pkg/shutdown"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// @title TempDB Server API
// @version 2.0
// @description Create temporary databases for testing and development
// @contact.name Thalles Passos
// @license.name MIT
// @BasePath /
func main() {
	// setup exit code for graceful shutdown
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	// load config
	env, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	// run the server
	cleanup, err := run(env)

	// run the cleanup after the server is terminated
	defer cleanup()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	// ensure the server is shutdown gracefully & app runs
	shutdown.Gracefully()
}

func run(env config.EnvVars) (func(), error) {
	app, cleanup, err := buildServer(env)
	if err != nil {
		return nil, err
	}

	// start the server
	go func() {
		app.Listen("0.0.0.0:" + env.PORT)
	}()

	// return a function to close the server and database
	return func() {
		cleanup()
		app.Shutdown()
	}, nil
}

func buildServer(env config.EnvVars) (*fiber.App, func(), error) {
	// init the storage
	db, err := storage.BootstrapPostgres(env.DATABASE_URI, 10*time.Second)
	if err != nil {
		return nil, nil, err
	}

	// create the fiber app
	app := fiber.New()

	// add middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// add health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy!")
	})

	// add docs
	app.Get("/swagger/*", swagger.HandlerDefault)

	// create the user domain
	databaseStorage := databases.NewDatabasesStorage(db)
	databaseController := databases.NewDatabasesController(databaseStorage, env.DATABASE_HOST)
	databases.AddTodoRoutes(app, databaseController)

	return app, func() {
		storage.ClosePostgres(db)
	}, nil
}
