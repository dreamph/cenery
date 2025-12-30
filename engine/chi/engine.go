package chi

import (
	"github.com/dreamph/cenery"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewApp() cenery.App {
	chiApp := chi.NewRouter()
	chiApp.Use(middleware.Recoverer)
	return New(chiApp)
}
