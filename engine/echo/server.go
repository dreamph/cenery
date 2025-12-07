package echo

import (
	"context"

	"github.com/dreamph/cenery"
	"github.com/labstack/echo/v4"
)

type app struct {
	server *echo.Echo
}

func New(server *echo.Echo) cenery.App {
	return &app{server: server}
}

func (a *app) Listen(addr string) error {
	return a.server.Start(addr)
}

func (a *app) Use(handlers ...cenery.Handler) {
	for _, handler := range handlers {
		h := handler // Copy variable to avoid closure capture bug
		a.server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				return a.processHandler(c, h, next)
			}
		})
	}
}

func (a *app) Get(path string, handlers ...cenery.Handler) {
	handler, middlewareHandlers := a.toHandlers(handlers)
	a.server.GET(path, handler, middlewareHandlers...)
}

func (a *app) Post(path string, handlers ...cenery.Handler) {
	handler, middlewareHandlers := a.toHandlers(handlers)
	a.server.POST(path, handler, middlewareHandlers...)
}

func (a *app) Put(path string, handlers ...cenery.Handler) {
	handler, middlewareHandlers := a.toHandlers(handlers)
	a.server.PUT(path, handler, middlewareHandlers...)
}

func (a *app) Delete(path string, handlers ...cenery.Handler) {
	handler, middlewareHandlers := a.toHandlers(handlers)
	a.server.DELETE(path, handler, middlewareHandlers...)
}

func (a *app) Head(path string, handlers ...cenery.Handler) {
	handler, middlewareHandlers := a.toHandlers(handlers)
	a.server.HEAD(path, handler, middlewareHandlers...)
}

func (a *app) Options(path string, handlers ...cenery.Handler) {
	handler, middlewareHandlers := a.toHandlers(handlers)
	a.server.OPTIONS(path, handler, middlewareHandlers...)
}

func (a *app) Connect(path string, handlers ...cenery.Handler) {
	handler, middlewareHandlers := a.toHandlers(handlers)
	a.server.CONNECT(path, handler, middlewareHandlers...)
}

func (a *app) Patch(path string, handlers ...cenery.Handler) {
	handler, middlewareHandlers := a.toHandlers(handlers)
	a.server.PATCH(path, handler, middlewareHandlers...)
}

func (a *app) Trace(path string, handlers ...cenery.Handler) {
	handler, middlewareHandlers := a.toHandlers(handlers)
	a.server.TRACE(path, handler, middlewareHandlers...)
}

func (a *app) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func (a *app) toHandlers(handlers []cenery.Handler) (func(c echo.Context) error, []echo.MiddlewareFunc) {
	var handler func(c echo.Context) error
	var middlewareHandlers []echo.MiddlewareFunc
	if len(handlers) == 1 {
		h := handlers[0]
		handler = func(c echo.Context) error {
			return a.processHandler(c, h, nil)
		}
	} else {
		h := handlers[len(handlers)-1]
		handler = func(c echo.Context) error {
			return a.processHandler(c, h, nil)
		}

		middlewares := handlers[:len(handlers)-1]
		for _, middleware := range middlewares {
			m := middleware // Copy variable to avoid closure capture bug
			middlewareHandler := func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					return a.processHandler(c, m, next)
				}
			}

			middlewareHandlers = append(middlewareHandlers, middlewareHandler)
		}
	}
	return handler, middlewareHandlers
}

func (a *app) processHandler(c echo.Context, handler cenery.Handler, next echo.HandlerFunc) error {
	svc := NewServerCtx(c, next)
	return handler(svc)
}
