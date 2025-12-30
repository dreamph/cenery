package main

import (
	"log"

	"github.com/dreamph/cenery"
	fasthttpengine "github.com/dreamph/cenery/fasthttp"
	"github.com/fasthttp/router"
)

func NewApp() cenery.App {
	routerApp := router.New()
	return fasthttpengine.New(routerApp)
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
