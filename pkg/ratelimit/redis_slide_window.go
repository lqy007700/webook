package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var slideWindowLua string

type RedisSlideWindow struct {
	cmd redis.Cmdable

	// 窗口大小
	interval time.Duration

	// 阈值
	rate int32

	// interval 内允许 rate 个请求
}

func NewRedisSlideWindow(cmd redis.Cmdable, interval time.Duration, rate int32) Limiter {
	return &RedisSlideWindow{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (r *RedisSlideWindow) Limited(ctx context.Context, key string) (bool, error) {
	r.cmd.Eval(ctx, slideWindowLua, []string{key}, 1, 1, 1)
	return true, nil
}
