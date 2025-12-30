package main

import (
	"fmt"
	"log"

	"github.com/dreamph/cenery"
	chiengine "github.com/dreamph/cenery/chi"
	echoengine "github.com/dreamph/cenery/echo"
	fasthttpengine "github.com/dreamph/cenery/fasthttp"
	fiberengine "github.com/dreamph/cenery/fiber"
	ginengine "github.com/dreamph/cenery/gin"
	"github.com/fasthttp/router"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
		JSONDecoder:       gojson.Unmarshal,
		JSONEncoder:       gojson.Marshal,
		StreamRequestBody: true,
	})
	fiberApp.Use(fiberrecover.New())
	return fiberengine.New(fiberApp)
}

func NewFiberApp() cenery.App {
	return fiberengine.NewApp()
}

func NewGinWithCustomize() cenery.App {
	ginApp := gin.New()
	ginApp.Use(gin.Recovery())
	return ginengine.New(ginApp)
}

func NewGinApp() cenery.App {
	return ginengine.NewApp()
}

func NewChiWithCustomize() cenery.App {
	chiApp := chi.NewRouter()
	chiApp.Use(middleware.Recoverer)
	return chiengine.New(chiApp)
}

func NewChiApp() cenery.App {
	return chiengine.NewApp()
}

func NewFastHTTPWithCustomize() cenery.App {
	routerApp := router.New()
	return fasthttpengine.New(routerApp)
}

func NewFastHTTPApp() cenery.App {
	return fasthttpengine.NewApp()
}

func main() {
	// for fiber
	//app := NewFiberApp() // or NewFiberAppWithCustomize()

	// for echo
	app := NewEchoApp() // or NewEchoWithCustomize()

	// for gin
	//app := NewGinApp() // or NewGinWithCustomize()

	// for chi
	//app := NewChiApp() // or NewChiWithCustomize()

	// for fasthttp
	//app := NewFastHTTPApp() // or NewFastHTTPWithCustomize()

	app.Use(func(c cenery.Ctx) error {
		//fmt.Println("global middleware")
		return c.Next()
	})

	app.Get("/", func(c cenery.Ctx) error {
		return c.SendString(200, "hello")
	})

	app.Get("/json", func(c cenery.Ctx) error {
		return c.SendJSON(200, "hello")
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
