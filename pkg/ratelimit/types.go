package ratelimit

import "context"

type Limiter interface {
	// Limited 是否出发限流
	// true 限流
	// err 错误
	Limited(ctx context.Context, key string) (bool, error)
}
