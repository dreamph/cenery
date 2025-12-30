package gin

import (
	"github.com/dreamph/cenery"
	"github.com/gin-gonic/gin"
)

func NewApp() cenery.App {
	ginApp := gin.New()
	ginApp.Use(gin.Recovery())
	return New(ginApp)
}
