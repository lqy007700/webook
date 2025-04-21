package repository

import (
	"context"
	"errors"
	"webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany        = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz string, code string, phone string) error
}

type CacheCodeRepository struct {
	cache cache.CodaCache
}

func NewCodeRepository(cache cache.CodaCache) CodeRepository {
	return &CacheCodeRepository{
		cache: cache,
	}
}

func (c *CacheCodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	err := c.cache.Set(ctx, biz, phone, code)
	switch {
	case errors.Is(err, ErrCodeSendTooMany):
		return ErrCodeSendTooMany
	default:
		return err
	}
}

// Verify 验证验证码的逻辑
func (c *CacheCodeRepository) Verify(ctx context.Context, biz string, code string, phone string) error {
	err := c.cache.Verify(ctx, biz, code, phone)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, cache.ErrCodeVerifyTooManyTimes):
		return ErrCodeVerifyTooManyTimes
	default:
		return err
	}
}
