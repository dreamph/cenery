package fiber

import (
	"context"
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
	return a.server.Listen(addr)
}

func (a *app) Use(handlers ...cenery.Handler) {
	apiHandlers := a.toHandlers(handlers...)
	var middlewareHandlers []interface{}
	for _, handler := range apiHandlers {
		middlewareHandlers = append(middlewareHandlers, handler)
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

func (a *app) processHandlers(c *fiber.Ctx, handlers ...cenery.Handler) error {
	svc := NewServerCtx(c)
	for _, handler := range handlers {
		err := handler(svc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *app) toHandlers(handlers ...cenery.Handler) []fiber.Handler {
	var handlerList []fiber.Handler
	for _, handler := range handlers {
		h := func(c *fiber.Ctx) error {
			svc := NewServerCtx(c)
			err := handler(svc)
			if err != nil {
				return err
			}
			return nil
		}
		handlerList = append(handlerList, h)
	}
	return handlerList
}
