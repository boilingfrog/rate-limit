package rateLimit

import (
	"context"
	"fmt"
	"math/rand"
	"rate-limit/redis"
	"time"
)

type LimitClient struct {
	rateLimit *redis.Redis
}

func New(conf *redis.Config) *LimitClient {
	return &LimitClient{
		rateLimit: redis.New(conf),
	}
}

type RateLimit interface {
	RateLimiter(ctx context.Context, param ...Param) error
	DefaultLimiter(ctx context.Context) error
	UserTimesLimiter(ctx context.Context, accountId string, maxThreads, expireTime int64) error
	TimesLimiter(ctx context.Context, key string, maxThreads, expireTime int64) error
	UserSingleRequestLimiter(ctx context.Context, accountId string, expireTime int64) error
}

func (p *LimitClient) RateLimiter(ctx context.Context, param ...Param) error {
	ps := evaluateParam(param)

	validAndAssignInput(ctx, ps)
	rand.Seed(time.Now().UnixNano())
	// 固定窗口的实现
	//res, err := p.rateLimit.Eval(ctx, fixedWindowScript(), []string{ps.Key}, ps.ExpireTime, ps.MaxThreads).Result()

	// 滑动窗口的实现
	res, err := p.rateLimit.Eval(ctx, slidingWindowScript(), []string{ps.Key}, ps.ExpireTime, ps.MaxThreads, time.Now().Unix(), rand.Int()).Result()
	if err != nil {
		return Err
	}
	if res.(int64) != 1 {
		return Err
	}

	if !ps.IsLimitTime {
		p.rateLimit.Del(ctx, ps.Key)
	}

	return nil
}

func (p *LimitClient) DefaultLimiter(ctx context.Context) error {
	return p.RateLimiter(ctx, nil)
}

func (p *LimitClient) UserTimesLimiter(ctx context.Context, accountId string, maxThreads, expireTime int64) error {
	return p.RateLimiter(ctx,
		MaxThreads(maxThreads),
		ExpireTime(expireTime),
		Key("User:limit:"+accountId),
		IsLimitUser(true),
		IsLimitTime(true),
	)
}

func (p *LimitClient) TimesLimiter(ctx context.Context, key string, maxThreads, expireTime int64) error {
	return p.RateLimiter(ctx,
		MaxThreads(maxThreads),
		ExpireTime(expireTime),
		Key(key),
		IsLimitTime(true),
	)
}

func (p *LimitClient) UserSingleRequestLimiter(ctx context.Context, accountId string, expireTime int64) error {
	// todo 业务使用的，可以根据用户id，加上ip,接口名来拼接 key
	return p.RateLimiter(ctx,
		MaxThreads(1),
		ExpireTime(expireTime),
		// 业务中使用需要完善
		Key(accountId),
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
