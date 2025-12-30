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
	app := cenery.NewServer(echoengine.NewApp())

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
app := cenery.NewServer(echoengine.NewApp())

// Fiber
app := cenery.NewServer(fiberengine.NewApp())

// Gin
app := cenery.NewServer(ginengine.NewApp())

// Chi
app := cenery.NewServer(chiengine.NewApp())

// fasthttp
app := cenery.NewServer(fasthttpengine.NewApp())
```

## Custom setup (when you need to tune)
```go
// Echo
echoApp := echo.New()
echoApp.Use(echomiddleware.Recover())
app := cenery.NewServer(echoengine.New(echoApp))

// Fiber
fiberApp := fiber.New(fiber.Config{
	JSONDecoder: gojson.Unmarshal,
	JSONEncoder: gojson.Marshal,
})
fiberApp.Use(fiberrecover.New())
app := cenery.NewServer(fiberengine.New(fiberApp))

// Gin
ginApp := gin.New()
ginApp.Use(gin.Recovery())
app := cenery.NewServer(ginengine.New(ginApp))

// Chi
chiApp := chi.NewRouter()
chiApp.Use(middleware.Recoverer)
app := cenery.NewServer(chiengine.New(chiApp))

// fasthttp
routerApp := router.New()
app := cenery.NewServer(fasthttpengine.New(routerApp))
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

## Route params
cenery accepts `:id` style and maps it for engines like Chi/fasthttp.
```go
app.Get("/users/:id", func(c cenery.Ctx) error {
	id := c.Params("id")
	return c.SendString(200, "id="+id)
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
