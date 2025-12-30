package chi

import (
	"context"
	"net/http"
	"os"

	"github.com/dreamph/cenery"
	"github.com/go-chi/chi/v5"
)

type app struct {
	server     *chi.Mux
	httpServer *http.Server
}

func New(server *chi.Mux) cenery.App {
	return &app{server: server}
}

func (a *app) Listen(addr string) error {
	_ = cenery.PrintLogo(os.Stdout)
	a.httpServer = &http.Server{
		Addr:    addr,
		Handler: a.server,
	}
	return a.httpServer.ListenAndServe()
}

func (a *app) Use(handlers ...cenery.Handler) {
	middlewares := a.toMiddlewares(handlers...)
	a.server.Use(middlewares...)
}

func (a *app) Get(path string, handlers ...cenery.Handler) {
	handler, middlewares := a.toHandlers(handlers...)
	a.server.With(middlewares...).Get(normalizePath(path), handler)
}

func (a *app) Post(path string, handlers ...cenery.Handler) {
	handler, middlewares := a.toHandlers(handlers...)
	a.server.With(middlewares...).Post(normalizePath(path), handler)
}

func (a *app) Put(path string, handlers ...cenery.Handler) {
	handler, middlewares := a.toHandlers(handlers...)
	a.server.With(middlewares...).Put(normalizePath(path), handler)
}

func (a *app) Delete(path string, handlers ...cenery.Handler) {
	handler, middlewares := a.toHandlers(handlers...)
	a.server.With(middlewares...).Delete(normalizePath(path), handler)
}

func (a *app) Head(path string, handlers ...cenery.Handler) {
	handler, middlewares := a.toHandlers(handlers...)
	a.server.With(middlewares...).Head(normalizePath(path), handler)
}

func (a *app) Options(path string, handlers ...cenery.Handler) {
	handler, middlewares := a.toHandlers(handlers...)
	a.server.With(middlewares...).Options(normalizePath(path), handler)
}

func (a *app) Connect(path string, handlers ...cenery.Handler) {
	handler, middlewares := a.toHandlers(handlers...)
	a.server.With(middlewares...).MethodFunc(http.MethodConnect, normalizePath(path), handler)
}

func (a *app) Patch(path string, handlers ...cenery.Handler) {
	handler, middlewares := a.toHandlers(handlers...)
	a.server.With(middlewares...).Patch(normalizePath(path), handler)
}

func (a *app) Trace(path string, handlers ...cenery.Handler) {
	handler, middlewares := a.toHandlers(handlers...)
	a.server.With(middlewares...).MethodFunc(http.MethodTrace, normalizePath(path), handler)
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

func (a *app) Shutdown(ctx context.Context) error {
	if a.httpServer == nil {
		return nil
	}
	return a.httpServer.Shutdown(ctx)
}

func (a *app) toHandlers(handlers ...cenery.Handler) (http.HandlerFunc, []func(http.Handler) http.Handler) {
	if len(handlers) == 0 {
		return func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}, nil
	}

	h := handlers[len(handlers)-1]
	handler := func(w http.ResponseWriter, r *http.Request) {
		svc := NewServerCtx(w, r, nil)
		if err := h(svc); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	middlewares := a.toMiddlewares(handlers[:len(handlers)-1]...)
	return handler, middlewares
}

func (a *app) toMiddlewares(handlers ...cenery.Handler) []func(http.Handler) http.Handler {
	if len(handlers) == 0 {
		return nil
	}

	middlewareHandlers := make([]func(http.Handler) http.Handler, len(handlers))
	for i, handler := range handlers {
		h := handler // Copy variable to avoid closure capture bug
		middlewareHandlers[i] = func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				svc := NewServerCtx(w, r, next)
				if err := h(svc); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			})
		}
	}
	return middlewareHandlers
}
