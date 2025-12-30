package fiber

import (
	"context"
	"os"

	"github.com/dreamph/cenery"
	"github.com/gofiber/fiber/v2"
)

type app struct {
	server *fiber.App
}

func New(server *fiber.App) cenery.App {
	return &app{server: server}
}

func (a *app) Listen(addr string) error {
	_ = cenery.PrintLogo(os.Stdout)
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
	a.server.Get(path, a.toHandlers(handlers...)...)
}

func (a *app) Post(path string, handlers ...cenery.Handler) {
	a.server.Post(path, a.toHandlers(handlers...)...)
}

func (a *app) Put(path string, handlers ...cenery.Handler) {
	a.server.Put(path, a.toHandlers(handlers...)...)
}

func (a *app) Delete(path string, handlers ...cenery.Handler) {
	a.server.Delete(path, a.toHandlers(handlers...)...)
}

func (a *app) Head(path string, handlers ...cenery.Handler) {
	a.server.Head(path, a.toHandlers(handlers...)...)
}

func (a *app) Options(path string, handlers ...cenery.Handler) {
	a.server.Options(path, a.toHandlers(handlers...)...)
}

func (a *app) Connect(path string, handlers ...cenery.Handler) {
	a.server.Connect(path, a.toHandlers(handlers...)...)
}

func (a *app) Patch(path string, handlers ...cenery.Handler) {
	a.server.Patch(path, a.toHandlers(handlers...)...)
}

func (a *app) Trace(path string, handlers ...cenery.Handler) {
	a.server.Trace(path, a.toHandlers(handlers...)...)
}

func (a *app) Shutdown(_ context.Context) error {
	return a.server.Shutdown()
}

func (a *app) toHandlers(handlers ...cenery.Handler) []fiber.Handler {
	handlerList := make([]fiber.Handler, len(handlers))
	for i, handler := range handlers {
		h := handler // Copy variable to avoid closure capture bug
		handlerList[i] = func(c *fiber.Ctx) error {
			svc := NewServerCtx(c)
			return h(svc)
		}
	}
	return handlerList
}
