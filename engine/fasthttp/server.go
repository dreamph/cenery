package fasthttp

import (
	"context"

	"github.com/dreamph/cenery"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type app struct {
	router      *router.Router
	server      *fasthttp.Server
	middlewares []cenery.Handler
}

func New(routerApp *router.Router) cenery.App {
	return &app{router: routerApp}
}

func (a *app) Name() string {
	return "fasthttp"
}

func (a *app) Listen(addr string) error {
	a.server = &fasthttp.Server{
		Handler: a.router.Handler,
	}
	return a.server.ListenAndServe(addr)
}

func (a *app) Use(handlers ...cenery.Handler) {
	a.middlewares = append(a.middlewares, handlers...)
}

func (a *app) Get(path string, handlers ...cenery.Handler) {
	a.router.GET(normalizePath(path), a.wrapWithMiddlewares(handlers...))
}

func (a *app) Post(path string, handlers ...cenery.Handler) {
	a.router.POST(normalizePath(path), a.wrapWithMiddlewares(handlers...))
}

func (a *app) Put(path string, handlers ...cenery.Handler) {
	a.router.PUT(normalizePath(path), a.wrapWithMiddlewares(handlers...))
}

func (a *app) Delete(path string, handlers ...cenery.Handler) {
	a.router.DELETE(normalizePath(path), a.wrapWithMiddlewares(handlers...))
}

func (a *app) Head(path string, handlers ...cenery.Handler) {
	a.router.HEAD(normalizePath(path), a.wrapWithMiddlewares(handlers...))
}

func (a *app) Options(path string, handlers ...cenery.Handler) {
	a.router.OPTIONS(normalizePath(path), a.wrapWithMiddlewares(handlers...))
}

func (a *app) Connect(path string, handlers ...cenery.Handler) {
	a.router.Handle(fasthttp.MethodConnect, normalizePath(path), a.wrapWithMiddlewares(handlers...))
}

func (a *app) Patch(path string, handlers ...cenery.Handler) {
	a.router.PATCH(normalizePath(path), a.wrapWithMiddlewares(handlers...))
}

func (a *app) Trace(path string, handlers ...cenery.Handler) {
	a.router.Handle(fasthttp.MethodTrace, normalizePath(path), a.wrapWithMiddlewares(handlers...))
}

func (a *app) Shutdown(ctx context.Context) error {
	if a.server == nil {
		return nil
	}
	return a.server.ShutdownWithContext(ctx)
}

func (a *app) wrapWithMiddlewares(handlers ...cenery.Handler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		all := make([]cenery.Handler, 0, len(a.middlewares)+len(handlers))
		all = append(all, a.middlewares...)
		all = append(all, handlers...)
		a.processHandlers(ctx, all, 0)
	}
}

func (a *app) processHandlers(ctx *fasthttp.RequestCtx, handlers []cenery.Handler, index int) {
	if index >= len(handlers) {
		return
	}

	handler := handlers[index]
	next := fasthttp.RequestHandler(func(c *fasthttp.RequestCtx) {
		a.processHandlers(c, handlers, index+1)
	})

	svc := NewServerCtx(ctx, next)
	if err := handler(svc); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

func (a *app) Handler() fasthttp.RequestHandler {
	return a.router.Handler
}

func normalizePath(path string) string {
	if path == "" {
		return path
	}

	var out []byte
	for i := 0; i < len(path); i++ {
		if path[i] != ':' {
			out = append(out, path[i])
			continue
		}

		if i > 0 && path[i-1] != '/' {
			out = append(out, path[i])
			continue
		}

		j := i + 1
		for j < len(path) && path[j] != '/' {
			j++
		}
		if j == i+1 {
			out = append(out, path[i])
			continue
		}

		out = append(out, '{')
		out = append(out, path[i+1:j]...)
		out = append(out, '}')
		i = j - 1
	}

	return string(out)
}
