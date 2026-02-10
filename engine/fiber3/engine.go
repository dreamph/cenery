package fiber3

import (
	"github.com/dreamph/cenery"
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	fiberrecover "github.com/gofiber/fiber/v3/middleware/recover"
)

func NewApp() cenery.App {
	fiberApp := fiber.New(fiber.Config{
		JSONDecoder: gojson.Unmarshal,
		JSONEncoder: gojson.Marshal,
	})
	fiberApp.Use(fiberrecover.New())
	return New(fiberApp)
}
