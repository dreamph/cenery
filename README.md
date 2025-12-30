# cenery

Switch engines. Keep your handlers. Ship fast.

cenery wraps popular Go web frameworks behind one clean API, so you can move
between Echo, Fiber, Gin, Chi, and fasthttp without rewriting your app.

## What you get
- One handler interface across engines
- Same middleware flow, different runtime
- Built-in helpers: JSON, streams, file uploads

## Install
```sh
go get github.com/dreamph/cenery
```

## Start in 20 lines
```go
package main

import (
	"log"

	"github.com/dreamph/cenery"
	echoengine "github.com/dreamph/cenery/echo"
)

func main() {
	app := echoengine.NewApp()

	app.Get("/", func(c cenery.Ctx) error {
		return c.SendString(200, "hello")
	})

	if err := app.Listen(":2000"); err != nil {
		log.Fatal(err)
	}
}
```

## Switch engines like a pro
```go
// Echo
app := echoengine.NewApp()

// Fiber
app := fiberengine.NewApp()

// Gin
app := ginengine.NewApp()

// Chi
app := chiengine.NewApp()

// fasthttp
app := fasthttpengine.NewApp()
```

## Custom setup (when you need to tune)
```go
// Echo
echoApp := echo.New()
echoApp.Use(echomiddleware.Recover())
app := echoengine.New(echoApp)

// Fiber
fiberApp := fiber.New(fiber.Config{
	JSONDecoder: gojson.Unmarshal,
	JSONEncoder: gojson.Marshal,
})
fiberApp.Use(fiberrecover.New())
app := fiberengine.New(fiberApp)

// Gin
ginApp := gin.New()
ginApp.Use(gin.Recovery())
app := ginengine.New(ginApp)

// Chi
chiApp := chi.NewRouter()
chiApp.Use(middleware.Recoverer)
app := chiengine.New(chiApp)

// fasthttp
routerApp := router.New()
app := fasthttpengine.New(routerApp)
```

## Helpers you will use a lot
```go
app.Post("/upload", func(c cenery.Ctx) error {
	file := c.FormFile("file")
	if file == nil {
		return c.SendJSON(400, map[string]string{"error": "missing file"})
	}
	return c.SendJSON(200, map[string]any{
		"name": file.FileName,
		"size": file.FileSize,
	})
})

app.Get("/stream", func(c cenery.Ctx) error {
	reader := strings.NewReader("streaming data")
	return c.SendStream(200, "text/plain", reader)
})
```

## Middleware example
```go
app.Use(func(c cenery.Ctx) error {
	start := time.Now()
	err := c.Next()
	elapsed := time.Since(start)
	log.Printf("method=%s path=%s status=%d in=%s",
		c.Request().GetHeader("Method"),
		c.Request().GetHeader("Path"),
		c.Response().GetHeader("Status"),
		elapsed,
	)
	return err
})
```

## Examples
Try these:
- `test/main.go`
- `test/cenery/echo/main.go`
- `test/cenery/fiber/main.go`
- `test/cenery/gin/main.go`
- `test/cenery/chi/main.go`
- `test/cenery/fasthttp/main.go`

## License
MIT
