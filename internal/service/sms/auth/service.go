package auth

import (
	"context"
	"webook/internal/service/sms"
)

type SMSService struct {
	svc sms.Service
}

func (S *SMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	panic("implement me")
}
