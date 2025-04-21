package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	"webook/internal/web"
)

func InitGin(handlerFunc []gin.HandlerFunc, handler *web.UserHandler) *gin.Engine {
	engine := gin.Default()
	handler.RegisterRouters(engine)
	engine.Use(handlerFunc...)
	return engine
}

func InitHandlerFunc() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowOrigins:     []string{"http://qie.com"},
			AllowMethods:     []string{"PUT"},
			AllowHeaders:     []string{"Content-Type", "authorization"},
			ExposeHeaders:    []string{"x-jwt-token"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				if strings.Contains(origin, "qie.com") {
					return true
				}
				return strings.Contains(origin, "localhost")
			},
			MaxAge: 12 * time.Hour,
		}),
	}
}
