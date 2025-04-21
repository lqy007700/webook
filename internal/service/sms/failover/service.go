package failover

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"webook/internal/service/sms"
)

type FailoverSMSService struct {
	svcs []sms.Service
	idx  uint64
}

func NewFailoverSMSService(svcs ...sms.Service) sms.Service {
	return &FailoverSMSService{
		svcs: svcs,
	}
}
func (f *FailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
		log.Printf("短信服务商%T调用失败，错误信息：%v", svc, err)
	}
	return errors.New("所有服务商全部失败")
}

func (f *FailoverSMSService) SendV1(ctx context.Context, tpl string, args []string, numbers ...string) error {
	ids := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := ids; i < ids+length; i++ {
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch {
		case err == nil:
			return nil
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return err
		}
		log.Printf("短信服务商%T调用失败，错误信息：%v", svc, err)
	}
	return errors.New("所有服务商全部失败")
}
