package echo

import (
	"github.com/dreamph/cenery"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func NewApp() cenery.App {
	echoApp := echo.New()
	echoApp.Use(echomiddleware.Recover())
	return New(echoApp)
}
