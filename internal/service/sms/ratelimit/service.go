package ratelimit

import (
	"context"
	"errors"
	"log"
	"webook/internal/service/sms"
	"webook/pkg/ratelimit"
)

type RateLimiterSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RateLimiterSMSService{svc: svc, limiter: limiter}
}

func (s *RateLimiterSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	ok, err := s.limiter.Limited(ctx, "sms")
	if err != nil {
		log.Printf("短信限流失败%v", err)
	}
	if ok {
		return errors.New("请求频繁")
	}
	return s.svc.Send(ctx, tpl, args, numbers...)
}
