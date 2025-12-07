# cenery

Golang Web Framework
- Dynamic change engine between Fiber and Echo
- Simple & Easy


Install
=======
``` sh
go get github.com/dreamph/cenery
```

Example
=======
``` go
package main

import (
	"fmt"
	"log"

	"github.com/dreamph/cenery"
	echoengine "github.com/dreamph/cenery/echo"
	fiberengine "github.com/dreamph/cenery/fiber"
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
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

func NewEchoWithCustomize() cenery.App {
	echoApp := echo.New()
	echoApp.Use(echomiddleware.Recover())
	return echoengine.New(echoApp)
}

func NewEchoApp() cenery.App {
	return echoengine.NewApp()
}

func NewFiberAppWithCustomize() cenery.App {
	fiberApp := fiber.New(fiber.Config{
		JSONDecoder: gojson.Unmarshal,
		JSONEncoder: gojson.Marshal,
	})
	fiberApp.Use(fiberrecover.New())
	return fiberengine.New(fiberApp)
}

func NewFiberApp() cenery.App {
	return fiberengine.NewApp()
}

func main() {
	app := NewFiberApp()
	// or
	// app := NewEchoApp()

	app.Use(func(c cenery.Ctx) error {
		fmt.Println("global middleware")
		return c.Next()
	})

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
		fmt.Println("handler uploading..")
		request := &UploadRequest{}
		err := c.BodyParser(request)
		if err != nil {
			return c.SendJSON(400, &ErrorResponse{Message: err.Error()})
		}

		request.File = c.FormFile("file")

		fmt.Println("handler upload success")

		c.Response().SetHeader("TEST", "1")

		return c.SendJSON(200, &UploadResponse{Name: request.Name, FileName: request.File.FileName})
	})

	err := app.Listen(":2000")
	if err != nil {
		log.Fatal(err.Error())
	}
}

//curl -X POST http://localhost:2000/create -d '{"name":"cenery"}'
//curl -v -F name=cenery -F file=@test.txt http://localhost:2000/upload


```


Buy Me a Coffee
=======
[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/dreamph)