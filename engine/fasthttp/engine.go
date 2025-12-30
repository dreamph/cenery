package fasthttp

import (
	"github.com/dreamph/cenery"
	"github.com/fasthttp/router"
)

func NewApp() cenery.App {
	routerApp := router.New()
	return New(routerApp)
}
