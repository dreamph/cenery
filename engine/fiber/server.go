package fiber

import (
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

func (a *app) Get(path string, handlers ...cenery.Handler) {
	a.server.Get(path, func(c *fiber.Ctx) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) Post(path string, handlers ...cenery.Handler) {
	a.server.Post(path, func(c *fiber.Ctx) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) Put(path string, handlers ...cenery.Handler) {
	a.server.Put(path, func(c *fiber.Ctx) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) Delete(path string, handlers ...cenery.Handler) {
	a.server.Delete(path, func(c *fiber.Ctx) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) Head(path string, handlers ...cenery.Handler) {
	a.server.Head(path, func(c *fiber.Ctx) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) Options(path string, handlers ...cenery.Handler) {
	a.server.Options(path, func(c *fiber.Ctx) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) Connect(path string, handlers ...cenery.Handler) {
	a.server.Connect(path, func(c *fiber.Ctx) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) Patch(path string, handlers ...cenery.Handler) {
	a.server.Patch(path, func(c *fiber.Ctx) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) Trace(path string, handlers ...cenery.Handler) {
	a.server.Trace(path, func(c *fiber.Ctx) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
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
