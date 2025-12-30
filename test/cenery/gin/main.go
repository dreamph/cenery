package main

import (
	"log"

	"github.com/dreamph/cenery"
	ginengine "github.com/dreamph/cenery/gin"
	"github.com/gin-gonic/gin"
)

func NewServerApp() cenery.ServerApp {
	ginApp := gin.New()
	ginApp.Use(gin.Recovery())
	return cenery.NewServer(ginengine.New(ginApp))
}

func main() {
	app := NewServerApp()

	app.Get("/", func(c cenery.Ctx) error {
		return c.SendString(200, "hello")
	})

	app.Get("/:id", func(c cenery.Ctx) error {
		return c.SendString(200, c.Params("id"))
	})

	err := app.Listen(":2003")
	if err != nil {
		log.Fatal(err.Error())
	}
}

//curl -X POST http://localhost:2000/create -d '{"name":"cenery"}'
//curl -v -F name=cenery -F file=@test.txt http://localhost:2000/upload
