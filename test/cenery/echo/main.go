package main

import (
	"log"

	echoengine "github.com/dreamph/cenery/echo"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/dreamph/cenery"
)

func NewServerApp() cenery.ServerApp {
	echoApp := echo.New()
	echoApp.Use(echomiddleware.Recover())
	return cenery.NewServer(echoengine.New(echoApp))
}

func main() {
	app := NewServerApp()

	app.Get("/", func(c cenery.Ctx) error {
		return c.SendString(200, "hello")
	})

	app.Get("/:id", func(c cenery.Ctx) error {
		return c.SendString(200, c.Params("id"))
	})

	err := app.Listen(":2002")
	if err != nil {
		log.Fatal(err.Error())
	}
}

//curl -X POST http://localhost:2000/create -d '{"name":"cenery"}'
//curl -v -F name=cenery -F file=@test.txt http://localhost:2000/upload
