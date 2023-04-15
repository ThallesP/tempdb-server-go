package databases

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DatabasesController struct {
	databaseHost string
	storage *DatabasesStorage
}

func NewDatabasesController(storage *DatabasesStorage, databaseHost string) *DatabasesController {
	return &DatabasesController{storage: storage, databaseHost: databaseHost}
}

type createDatabaseInput struct {
	DatabaseType string `json:"database_type"`
	ExpiresInMilliseconds int `json:"expires_in_milliseconds"`
}

type createDatabaseResponse struct {
	Host string `json:"host"`
	DatabaseName string `json:"database_name"`
	User string `json:"user"`
	Password string `json:"password"`
	ConnectionString string `json:"connection_string"`
	ExpiresIn string `json:"expires_in"`
}

// @Summary Create temporary database
// @Description creates temporary dabase
// @Tags databases
// @Accept application/json
// @Produce json
// @Param database body createDatabaseInput true "Database details"
// @Success 200 {object} createDatabaseResponse
// @Router /databases [post]
func (d *DatabasesController) create(c *fiber.Ctx) error {
	var input createDatabaseInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if input.DatabaseType != "postgres" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Database type not supported",
		})
	}

	databaseName, err := d.storage.Create()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create database",
		})
	}

	credentials, err := d.storage.CreateUser(databaseName)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create database user",
			"error": err.Error(),
		})
	}

	go func() {
		<-time.After(time.Duration(input.ExpiresInMilliseconds) * time.Millisecond)
		fmt.Println("Deleting database", databaseName)
		d.storage.Delete(databaseName)
	}()

	return c.Status(fiber.StatusOK).JSON(createDatabaseResponse{
		Host: d.databaseHost,
		DatabaseName: databaseName,
		User: credentials.Username,
		Password: credentials.Password,
		ConnectionString: fmt.Sprintf("postgres://%s:%s@%s/%s", credentials.Username, credentials.Password, d.databaseHost, databaseName),
		ExpiresIn: time.Now().Add(time.Duration(input.ExpiresInMilliseconds)).Format(time.RFC1123),
	})
}