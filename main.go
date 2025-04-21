package main

import (
	"github.com/gin-gonic/gin"
	"webook/ioc"

	//"context"
	"github.com/redis/go-redis/v9"
	//"strings"
	//"time"
	//"webook/config"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/service/sms/memory"
	"webook/internal/web"

	//"github.com/gin-contrib/sessions"

	//"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ErrKeyNotExist = redis.Nil

func main() {
	s := InitWebserver()
	err := s.Run(":8081")
	if err != nil {
		panic(err)
	}
}

func InitWebserver() *gin.Engine {
	v := ioc.InitHandlerFunc()
	db := ioc.Initdb()
	gormUserDAO := dao.NewUserDAO(db)
	cmdable := ioc.InitRedis()
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(gormUserDAO, userCache)
	userService := service.NewUserService(userRepository)
	codaCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codaCache)
	smsService := ioc.NewSmsService()
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService)
	engine := ioc.InitGin(v, userHandler)
	return engine
}

// func initWebServer() *gin.Engine {
// gin engine
// s := gin.Default()

//s.Use(cors.New(cors.Config{
//	AllowOrigins:     []string{"http://qie.com"},
//	AllowMethods:     []string{"PUT"},
//	AllowHeaders:     []string{"Content-Type", "authorization"},
//	ExposeHeaders:    []string{"x-jwt-token"},
//	AllowCredentials: true,
//	AllowOriginFunc: func(origin string) bool {
//		if strings.Contains(origin, "qie.com") {
//			return true
//		}
//		return strings.Contains(origin, "localhost")
//	},
//	MaxAge: 12 * time.Hour,
//}))

//store, err := redis.NewStore(10, "tcp", config.Config.Redis.Addr, "", []byte("secret"))
//if err != nil {
//	panic(err)
//}
//s.Use(sessions.Sessions("sid", store))

// s.Use(middleware.NewLoginMiddlewareBuilder().
// IgnorePaths("/users/signup").IgnorePaths("/users/login").IgnorePaths("/users/loginjwt").
// Build())

//fmt.Println(1)
//
//s.Use(middleware.NewLoginJwtMiddlewareBuilder().
//	IgnorePaths("/users/loginjwt").
//	IgnorePaths("/users/signup").
//	Build())
// return s
// }

func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	userDao := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(rdb)
	repo := repository.NewUserRepository(userDao, userCache)
	svc := service.NewUserService(repo)

	smsSvc := &memory.Service{}

	codeCache := cache.NewCodeCache(rdb)
	codeRepo := repository.NewCodeRepository(codeCache)
	codeSvc := service.NewCodeService(codeRepo, smsSvc)
	return web.NewUserHandler(svc, codeSvc)
}
