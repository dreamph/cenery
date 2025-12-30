package main

import (
	"log"

	"github.com/dreamph/cenery"
	fasthttpengine "github.com/dreamph/cenery/fasthttp"
	"github.com/fasthttp/router"
)

func NewServerApp() cenery.ServerApp {
	routerApp := router.New()
	return cenery.NewServer(fasthttpengine.New(routerApp))
}

func main() {
	app := NewServerApp()

	app.Get("/", func(c cenery.Ctx) error {
		return c.SendString(200, "hello")
	})

	err := app.Listen(":2005")
	if err != nil {
		log.Fatal(err.Error())
	}
}

//curl -X POST http://localhost:2000/create -d '{"name":"cenery"}'
//curl -v -F name=cenery -F file=@test.txt http://localhost:2000/upload
