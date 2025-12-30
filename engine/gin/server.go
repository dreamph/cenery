package gin

import (
	"context"
	"net/http"

	"github.com/dreamph/cenery"
	"github.com/gin-gonic/gin"
)

type app struct {
	server     *gin.Engine
	httpServer *http.Server
}

func New(server *gin.Engine) cenery.App {
	return &app{server: server}
}

func (a *app) Name() string {
	return "Gin"
}

func (a *app) Listen(addr string) error {
	a.httpServer = &http.Server{
		Addr:    addr,
		Handler: a.server,
	}
	return a.httpServer.ListenAndServe()
}

func (a *app) Use(handlers ...cenery.Handler) {
	a.server.Use(a.toHandlers(handlers...)...)
}

func (a *app) Get(path string, handlers ...cenery.Handler) {
	a.server.GET(path, a.toHandlers(handlers...)...)
}

func (a *app) Post(path string, handlers ...cenery.Handler) {
	a.server.POST(path, a.toHandlers(handlers...)...)
}

func (a *app) Put(path string, handlers ...cenery.Handler) {
	a.server.PUT(path, a.toHandlers(handlers...)...)
}

func (a *app) Delete(path string, handlers ...cenery.Handler) {
	a.server.DELETE(path, a.toHandlers(handlers...)...)
}

func (a *app) Head(path string, handlers ...cenery.Handler) {
	a.server.HEAD(path, a.toHandlers(handlers...)...)
}

func (a *app) Options(path string, handlers ...cenery.Handler) {
	a.server.OPTIONS(path, a.toHandlers(handlers...)...)
}

func (a *app) Connect(path string, handlers ...cenery.Handler) {
	a.server.Handle(http.MethodConnect, path, a.toHandlers(handlers...)...)
}

func (a *app) Patch(path string, handlers ...cenery.Handler) {
	a.server.PATCH(path, a.toHandlers(handlers...)...)
}

func (a *app) Trace(path string, handlers ...cenery.Handler) {
	a.server.Handle(http.MethodTrace, path, a.toHandlers(handlers...)...)
}

func (a *app) Shutdown(ctx context.Context) error {
	if a.httpServer == nil {
		return nil
	}
	return a.httpServer.Shutdown(ctx)
}

func (a *app) toHandlers(handlers ...cenery.Handler) []gin.HandlerFunc {
	handlerList := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		h := handler // Copy variable to avoid closure capture bug
		handlerList[i] = func(c *gin.Context) {
			svc := NewServerCtx(c)
			if err := h(svc); err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}
	return handlerList
}
