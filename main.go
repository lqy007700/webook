package main

import (
	"net/http"
	"strings"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//db := initDB()
	//u := initUser(db)
	//s := initWebServer()
	//u.RegisterRouters(s)
	s := gin.Default()
	s.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello World")
	})
	s.Run(":8081")
}

func initWebServer() *gin.Engine {
	// gin engine
	s := gin.Default()

	s.Use(cors.New(cors.Config{
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
	}))

	store, err := redis.NewStore(10, "tcp", "localhost:16379", "", []byte("secret"))
	if err != nil {
		panic(err)
	}
	s.Use(sessions.Sessions("sid", store))

	// s.Use(middleware.NewLoginMiddlewareBuilder().
	// IgnorePaths("/users/signup").IgnorePaths("/users/login").IgnorePaths("/users/loginjwt").
	// Build())

	s.Use(middleware.NewLoginJwtMiddlewareBuilder().
		IgnorePaths("/users/loginjwt").
		Build())
	return s
}

func initUser(db *gorm.DB) *web.UserHandler {
	dao := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(dao)
	svc := service.NewUserService(repo)
	return web.NewUserHandler(svc)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13306)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
