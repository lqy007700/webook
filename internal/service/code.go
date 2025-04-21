package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"webook/internal/repository"
	"webook/internal/service/sms"
)

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany
var ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes

//type CodeService interface {
//	Send(ctx context.Context, biz, phone string) error
//	Verify(ctx context.Context, biz string, code string, phone string) error
//}

type CodeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	// 生成验证码
	// 存储 redis
	// 发送
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	// 发送
	return svc.smsSvc.Send(ctx, biz, []string{code}, phone)
}

func (svc *CodeService) generateCode() string {
	num := rand.Intn(999999)
	return fmt.Sprintf("%06d", num)
}

func (svc *CodeService) Verify(ctx context.Context, biz string, code string, phone string) error {
	err := svc.repo.Verify(ctx, biz, code, phone)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, repository.ErrCodeVerifyTooManyTimes):
		return ErrCodeVerifyTooManyTimes
	default:
		return err
	}
}
