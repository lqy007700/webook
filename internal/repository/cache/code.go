package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodaCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz string, code string, phone string) error
}

type RedisCodeCache struct {
	client redis.Cmdable
}

var (
	ErrCodeSendTooMany        = errors.New("send too many")
	ErrCodeVerifyTooManyTimes = errors.New("verify too many")
)

func NewCodeCache(client redis.Cmdable) CodaCache {
	return &RedisCodeCache{
		client: client,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}

	switch res {
	case 0:
		return nil
	case -1:
		return ErrCodeSendTooMany
	default:
		return fmt.Errorf("system err")
	}

}

func (c *RedisCodeCache) Verify(ctx context.Context, biz string, code string, phone string) error {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{key(biz, phone)}, code).Int()
	switch {
	case err != nil:
		return err
	case res == 0:
		return nil
	case res == -1:
		return ErrCodeVerifyTooManyTimes
	default:
		return fmt.Errorf("system err")
	}
}

func key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
