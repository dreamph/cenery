package main

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func main() {
	// Echo instance
	e := echo.New()

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Start server
	log.Fatal(e.Start(":3001"))
}
