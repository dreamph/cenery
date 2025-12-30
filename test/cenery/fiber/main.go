package main

import (
	"log"

	"github.com/dreamph/cenery"
	fiberengine "github.com/dreamph/cenery/fiber"
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
)

func NewServerApp() cenery.ServerApp {
	fiberApp := fiber.New(fiber.Config{
		JSONDecoder: gojson.Unmarshal,
		JSONEncoder: gojson.Marshal,
	})
	fiberApp.Use(fiberrecover.New())
	return cenery.NewServer(fiberengine.New(fiberApp))
}

func main() {
	app := NewServerApp()

	app.Get("/", func(c cenery.Ctx) error {
		return c.SendString(200, "hello")
	})

	app.Get("/:id", func(c cenery.Ctx) error {
		return c.SendString(200, c.Params("id"))
	})

	err := app.Listen(":2001")
	if err != nil {
		log.Fatal(err.Error())
	}
}

//curl -X POST http://localhost:2000/create -d '{"name":"cenery"}'
//curl -v -F name=cenery -F file=@test.txt http://localhost:2000/upload
