package main

import (
	"log"

	"github.com/dreamph/cenery"
	chiengine "github.com/dreamph/cenery/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewApp() cenery.App {
	chiApp := chi.NewRouter()
	chiApp.Use(middleware.Recoverer)
	return chiengine.New(chiApp)
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
