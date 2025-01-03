package main

import (
	echoengine "github.com/dreamph/cenery/engine/echo"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"log"

	"github.com/dreamph/cenery"
)

func NewApp() cenery.App {
	echoApp := echo.New()
	echoApp.Use(echomiddleware.Recover())
	return echoengine.New(echoApp)
}

func main() {
	app := NewApp()

	app.Get("/", func(c cenery.Ctx) error {
		return c.SendString(200, "hello")
	})

	err := app.Listen(":2000")
	if err != nil {
		log.Fatal(err.Error())
	}
}

//curl -X POST http://localhost:2000/create -d '{"name":"cenery"}'
//curl -v -F name=cenery -F file=@test.txt http://localhost:2000/upload
