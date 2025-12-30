package cenery

import (
	"context"
	"os"
)

type ServerApp interface {
	Get(path string, handlers ...Handler)
	Post(path string, handlers ...Handler)
	Put(path string, handlers ...Handler)
	Delete(path string, handlers ...Handler)
	Head(path string, handlers ...Handler)
	Options(path string, handlers ...Handler)
	Connect(path string, handlers ...Handler)
	Patch(path string, handlers ...Handler)
	Trace(path string, handlers ...Handler)
	Use(handlers ...Handler)
	Listen(addr string) error
	Shutdown(ctx context.Context) error
}

type serverApp struct {
	app App
}

func (s *serverApp) Connect(path string, handlers ...Handler) {
	s.app.Connect(path, handlers...)
}

func (s *serverApp) Delete(path string, handlers ...Handler) {
	s.app.Delete(path, handlers...)
}

func (s *serverApp) Get(path string, handlers ...Handler) {
	s.app.Get(path, handlers...)
}

func (s *serverApp) Head(path string, handlers ...Handler) {
	s.app.Head(path, handlers...)
}

func (s *serverApp) Options(path string, handlers ...Handler) {
	s.app.Options(path, handlers...)
}

func (s *serverApp) Patch(path string, handlers ...Handler) {
	s.app.Patch(path, handlers...)
}

func (s *serverApp) Post(path string, handlers ...Handler) {
	s.app.Post(path, handlers...)
}

func (s *serverApp) Put(path string, handlers ...Handler) {
	s.app.Put(path, handlers...)
}

func (s *serverApp) Shutdown(ctx context.Context) error {
	return s.app.Shutdown(ctx)
}

func (s *serverApp) Trace(path string, handlers ...Handler) {
	s.app.Trace(path, handlers...)
}

func (s *serverApp) Use(handlers ...Handler) {
	s.app.Use(handlers...)
}

func (s *serverApp) Listen(addr string) error {
	_ = PrintLogo(os.Stdout, s.app.Name(), addr)
	return s.app.Listen(addr)
}

func NewServer(app App) ServerApp {
	return &serverApp{
		app: app,
	}
}
