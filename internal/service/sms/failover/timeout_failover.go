package failover

import (
	"context"
	"errors"
	"sync/atomic"
	"webook/internal/service/sms"
)

type TimeoutFailoverSMSService struct {
	svcs []sms.Service
	// 连续超时的个数
	idx       int32
	cnt       int32
	threshold int32
}

func NewTimeoutFailoverSMSService() sms.Service {
	return &TimeoutFailoverSMSService{
		svcs: []sms.Service{},
		cnt:  10,
	}
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for {
		idx := atomic.LoadInt32(&t.cnt)
		cnt := atomic.LoadInt32(&t.cnt)
		if cnt > t.threshold {
			newIdx := (idx + 1) % int32(len(t.svcs))
			if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
				atomic.StoreInt32(&t.cnt, 0)
				idx = newIdx
			} else {
				// 重试
				continue
			}
		}
		err := t.svcs[idx].Send(ctx, tpl, args, numbers...)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			atomic.AddInt32(&t.cnt, 1)
			return err
		case err == nil:
			atomic.StoreInt32(&t.cnt, 0)
			return nil
		}
		return err
	}
}
