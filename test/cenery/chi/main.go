package main

import (
	"log"

	"github.com/dreamph/cenery"
	chiengine "github.com/dreamph/cenery/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServerApp() cenery.ServerApp {
	chiApp := chi.NewRouter()
	chiApp.Use(middleware.Recoverer)
	return cenery.NewServer(chiengine.New(chiApp))
}

func main() {
	app := NewServerApp()

	app.Get("/", func(c cenery.Ctx) error {
		return c.SendString(200, "hello")
	})

	err := app.Listen(":2004")
	if err != nil {
		log.Fatal(err.Error())
	}
}

//curl -X POST http://localhost:2000/create -d '{"name":"cenery"}'
//curl -v -F name=cenery -F file=@test.txt http://localhost:2000/upload
