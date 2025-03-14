package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"webook/internal/web"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginJwtMiddlewareBuilder struct {
	paths []string
}

func NewLoginJwtMiddlewareBuilder() *LoginJwtMiddlewareBuilder {
	return &LoginJwtMiddlewareBuilder{}
}

func (l *LoginJwtMiddlewareBuilder) IgnorePaths(path string) *LoginJwtMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJwtMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if path == ctx.Request.URL.Path {
				return
			}
		}
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claim := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claim, func(t *jwt.Token) (interface{}, error) {
			return []byte("eHwX09d&*3KLs0^lm#PqA5RzVcT7NyU4QbFiGj2M8W!n@tYh"), nil
		})
		if err != nil || token == nil || !token.Valid || claim.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		fmt.Println(claim.UserAgent)
		fmt.Println(ctx.Request.UserAgent())
		if claim.UserAgent != ctx.Request.UserAgent() {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 每半个月刷新一次 token
		now := time.Now()
		if claim.ExpiresAt.Sub(now) < time.Hour*24*15 {
			claim.ExpiresAt = jwt.NewNumericDate(time.Now().AddDate(0, 1, 0))
			tokenStr, err = token.SignedString([]byte("eHwX09d&*3KLs0^lm#PqA5RzVcT7NyU4QbFiGj2M8W!n@tYh"))
			if err != nil {
				// 续约失败
			}
			ctx.Header("x-jwt-token", tokenStr)
		}

		ctx.Set("claim", claim)

	}
}
