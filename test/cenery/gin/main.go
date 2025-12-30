package main

import (
	"log"

	"github.com/dreamph/cenery"
	ginengine "github.com/dreamph/cenery/gin"
	"github.com/gin-gonic/gin"
)

func NewApp() cenery.App {
	ginApp := gin.New()
	ginApp.Use(gin.Recovery())
	return ginengine.New(ginApp)
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
