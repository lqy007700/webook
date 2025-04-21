package tencent

import (
	"context"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId     string
	signature string
	client    *sms.Client
}

func NewService(client *sms.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	return nil
}
