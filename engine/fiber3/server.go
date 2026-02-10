package fiber3

import (
	"context"

	"github.com/dreamph/cenery"
	"github.com/gofiber/fiber/v3"
)

type app struct {
	server *fiber.App
}

func New(server *fiber.App) cenery.App {
	return &app{server: server}
}

func (a *app) Name() string {
	return "Fiber v3"
}

func (a *app) Listen(addr string) error {
	return a.server.Listen(addr)
}

func (a *app) Use(handlers ...cenery.Handler) {
	apiHandlers := a.toHandlers(handlers...)
	middlewareHandlers := make([]any, len(apiHandlers))
	for i, handler := range apiHandlers {
		middlewareHandlers[i] = handler
	}
	a.server.Use(middlewareHandlers...)
}

func (a *app) Get(path string, handlers ...cenery.Handler) {
	first, rest := a.toRouteHandlers(handlers...)
	a.server.Get(path, first, rest...)
}

func (a *app) Post(path string, handlers ...cenery.Handler) {
	first, rest := a.toRouteHandlers(handlers...)
	a.server.Post(path, first, rest...)
}

func (a *app) Put(path string, handlers ...cenery.Handler) {
	first, rest := a.toRouteHandlers(handlers...)
	a.server.Put(path, first, rest...)
}

func (a *app) Delete(path string, handlers ...cenery.Handler) {
	first, rest := a.toRouteHandlers(handlers...)
	a.server.Delete(path, first, rest...)
}

func (a *app) Head(path string, handlers ...cenery.Handler) {
	first, rest := a.toRouteHandlers(handlers...)
	a.server.Head(path, first, rest...)
}

func (a *app) Options(path string, handlers ...cenery.Handler) {
	first, rest := a.toRouteHandlers(handlers...)
	a.server.Options(path, first, rest...)
}

func (a *app) Connect(path string, handlers ...cenery.Handler) {
	first, rest := a.toRouteHandlers(handlers...)
	a.server.Connect(path, first, rest...)
}

func (a *app) Patch(path string, handlers ...cenery.Handler) {
	first, rest := a.toRouteHandlers(handlers...)
	a.server.Patch(path, first, rest...)
}

func (a *app) Trace(path string, handlers ...cenery.Handler) {
	first, rest := a.toRouteHandlers(handlers...)
	a.server.Trace(path, first, rest...)
}

func (a *app) Shutdown(_ context.Context) error {
	return a.server.Shutdown()
}

func (a *app) toHandlers(handlers ...cenery.Handler) []fiber.Handler {
	handlerList := make([]fiber.Handler, len(handlers))
	for i, handler := range handlers {
		h := handler // Copy variable to avoid closure capture bug
		handlerList[i] = func(c fiber.Ctx) error {
			svc := NewServerCtx(c)
			return h(svc)
		}
	}
	return handlerList
}

func (a *app) toRouteHandlers(handlers ...cenery.Handler) (any, []any) {
	handlerList := a.toHandlers(handlers...)
	if len(handlerList) == 0 {
		return fiber.Handler(func(c fiber.Ctx) error { return c.Next() }), nil
	}

	first := any(handlerList[0])
	rest := make([]any, 0, len(handlerList)-1)
	for _, h := range handlerList[1:] {
		rest = append(rest, h)
	}
	return first, rest
}
