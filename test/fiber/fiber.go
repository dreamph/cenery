package main

import (
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
	"log"
)

func main() {
	app := fiber.New(fiber.Config{
		JSONDecoder: gojson.Unmarshal,
		JSONEncoder: gojson.Marshal,
	})
	app.Use(fiberrecover.New())

	// Define a route for the GET method on the root path '/'
	app.Get("/", func(c *fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
