package ioc

import (
	"webook/internal/service/sms"
	"webook/internal/service/sms/memory"
)

func NewSmsService() sms.Service {
	return memory.NewService()
}
