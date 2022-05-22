package rateLimit

import (
	"context"
	"fmt"
	"rate-limit/redis"
)

type LimitClient struct {
	rateLimit redis.Redis
}

func New(conf *redis.Config) *LimitClient {
	return &LimitClient{
		rateLimit: redis.New(conf),
	}
}

type RateLimit interface {
	RateLimiter(param ...Param) error
	DefaultLimiter() error
	UserTimesLimiter(accountId string, maxThreads, expireTime int64) error
	TimesLimiter(key string, maxThreads, expireTime int64) error
	SingleRequestLimiter(key string, expireTime int64) error
}

func (p *LimitClient) RateLimiter(param ...Param) error {
	ctx := context.Background()
	ps := evaluateParam(param)

	validAndAssignInput(ctx, ps)

	pipe := p.rateLimit.Pipeline(ctx)
	pipe.Send("INCR", ps.Key)
	pipe.Send("TTL", ps.Key)

	replies, err := pipe.Receive()
	if err != nil {
		return RateLimitErr
	}

	var (
		current = replies[0].(int64)
		ttl     = replies[1].(int64)
	)

	if current == int64(1) || ttl == int64(-1) {
		p.rateLimit.Do(ctx, "EXPIRE", ps.Key, ps.ExpireTime)
	}

	if current > ps.MaxThreads {
		return RateLimitErr
	}

	if !ps.IsLimitTime {
		p.rateLimit.Del(ctx, ps.Key)
	}

	return nil
}

func (p *LimitClient) DefaultLimiter() error {
	return p.RateLimiter(nil)
}

func (p *LimitClient) UserTimesLimiter(accountId string, maxThreads, expireTime int64) error {
	return p.RateLimiter(
		MaxThreads(maxThreads),
		ExpireTime(expireTime),
		Key("User:limit:"+accountId),
		IsLimitUser(true),
		IsLimitTime(true),
	)
}

func (p *LimitClient) TimesLimiter(key string, maxThreads, expireTime int64) error {
	return p.RateLimiter(
		MaxThreads(maxThreads),
		ExpireTime(expireTime),
		Key(key),
		IsLimitTime(true),
	)
}

func (p *LimitClient) SingleRequestLimiter(key string, expireTime int64) error {
	return p.RateLimiter(
		MaxThreads(1),
		ExpireTime(expireTime),
		Key(key),
		IsLimitTime(false),
	)
}

func validAndAssignInput(ctx context.Context, p *params) {
	keyItem := "common"

	if p.ExpireTime == 0 {
		p.ExpireTime = DefaultExpireTime
	}

	if p.MaxThreads == 0 {
		p.MaxThreads = DefaultMaxThreads
	}

	if p.Key == "" {
		// 可以根据路由的请求，例如接口名，ip,对默认的key进行定制
		p.Key = fmt.Sprintf("%s:%s", DefaultPrefix, keyItem)
	}
}
