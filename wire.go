//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/ioc"
)

func InitWebserver() *gin.Engine {
	wire.Build(
		ioc.Initdb,
		ioc.InitRedis,

		dao.NewUserDAO,
		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		service.NewUserService,
		service.NewCodeService,
		ioc.NewSmsService,

		web.NewUserHandler,
		ioc.InitGin,
		ioc.InitHandlerFunc,
	)
	return new(gin.Engine)
}
