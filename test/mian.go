package main

import (
	"github.com/dreamph/cenery"
	echoengine "github.com/dreamph/cenery/engine/echo"
	fiberengine "github.com/dreamph/cenery/engine/fiber"
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"log"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type CreateRequest struct {
	Name string `json:"name"`
}

type CreateResponse struct {
	Data string `json:"data"`
}

type UploadRequest struct {
	Name string           `json:"name"`
	File *cenery.FileData `json:"-"`
}

type UploadResponse struct {
	Name     string `json:"name"`
	FileName string `json:"fileName"`
}

func NewEcho() cenery.App {
	echoApp := echo.New()
	echoApp.Use(echomiddleware.Recover())
	return echoengine.New(echoApp)
}

func NewFiber() cenery.App {
	fiberApp := fiber.New(fiber.Config{
		JSONDecoder: gojson.Unmarshal,
		JSONEncoder: gojson.Marshal,
	})
	fiberApp.Use(fiberrecover.New())
	return fiberengine.New(fiberApp)
}

func main() {
	app := NewFiber() //NewEcho()

	app.Get("/", func(c cenery.Ctx) error {
		return c.SendString(200, "hello")
	})

	app.Post("/create", func(c cenery.Ctx) error {
		request := &CreateRequest{}
		err := c.BodyParser(request)
		if err != nil {
			return c.SendJSON(400, &ErrorResponse{Message: err.Error()})
		}
		return c.SendJSON(200, &CreateResponse{Data: "welcome : " + request.Name})
	})

	app.Post("/upload", func(c cenery.Ctx) error {
		request := &UploadRequest{}
		err := c.BodyParser(request)
		if err != nil {
			return c.SendJSON(400, &ErrorResponse{Message: err.Error()})
		}

		request.File = c.FormFile("file")

		return c.SendJSON(200, &UploadResponse{Name: request.Name, FileName: request.File.FileName})
	})

	err := app.Listen(":2000")
	if err != nil {
		log.Fatal(err.Error())
	}
}

//curl -X POST http://localhost:2000/create -d '{"name":"cenery"}'
//curl -v -F name=cenery -F file=@test.txt http://localhost:2000/upload
