package echo

import (
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

func (a *app) Get(path string, handlers ...cenery.Handler) {
	a.server.GET(path, func(c echo.Context) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) Post(path string, handlers ...cenery.Handler) {
	a.server.POST(path, func(c echo.Context) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) Put(path string, handlers ...cenery.Handler) {
	a.server.PUT(path, func(c echo.Context) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) Delete(path string, handlers ...cenery.Handler) {
	a.server.DELETE(path, func(c echo.Context) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) Head(path string, handlers ...cenery.Handler) {
	a.server.HEAD(path, func(c echo.Context) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) Options(path string, handlers ...cenery.Handler) {
	a.server.OPTIONS(path, func(c echo.Context) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) Connect(path string, handlers ...cenery.Handler) {
	a.server.CONNECT(path, func(c echo.Context) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) Patch(path string, handlers ...cenery.Handler) {
	a.server.PATCH(path, func(c echo.Context) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) Trace(path string, handlers ...cenery.Handler) {
	a.server.TRACE(path, func(c echo.Context) error {
		err := a.processHandlers(c, handlers...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *app) processHandlers(c echo.Context, handlers ...cenery.Handler) error {
	svc := NewServerCtx(c)
	for _, handler := range handlers {
		err := handler(svc)
		if err != nil {
			return err
		}
	}
	return nil
}
