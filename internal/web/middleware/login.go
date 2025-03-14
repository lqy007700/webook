package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if path == ctx.Request.URL.Path {
				return
			}
		}

		sess := sessions.Default(ctx)
		id := sess.Get("userId")

		if id == nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		sess.Options(sessions.Options{
			MaxAge: 60,
		})

		// 获取上一次的更新时间
		update_time := sess.Get("update_time")
		now := time.Now().UnixMilli()
		if update_time == nil {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
		update_time_val, _ := update_time.(int64)
		if now-update_time_val > 60*1000 {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
	}
}
