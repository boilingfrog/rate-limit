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

// 这是固定窗口的限流实现
// 栗子 对ip做限流 3秒内只能访问100次
// 固定窗口就是 第1秒访问了1次，第三秒访问了99次，那么第3秒仍然可访问100次
func fixedWindowScript() string {
	script := `
		local expireTime = ARGV[1] 
		local limitNum = ARGV[2]
		local key = KEYS[1]

		local visitNum = redis.call('incr', key)
		if visitNum == 1 then
				redis.call('expire', key, expireTime)
		end
		
		if visitNum > tonumber(limitNum) then
				return 0
		end
		
		return 1;
    `
	return script
}

// 滑动窗口的实现
func slidingWindowScript() string {
	script := `
		-- 过期时间，单位秒
		local expireTime = ARGV[1]
		-- 限制的最大数量
		local maxThreads = ARGV[2]
		-- 排序的分数
		local score = ARGV[3]
		-- 添加的成员
		local randomValue = ARGV[4]

		-- 有序集合的key
		local key = KEYS[1]
		-- 当前计算开始的分数
		local beginScore = tonumber(score)-tonumber(expireTime)

		local visitNum = redis.call('ZCOUNT', key, beginScore, tonumber(score))
		if visitNum >= tonumber(maxThreads) then
			return 0
		end

		redis.call('ZADD', key, score, randomValue)

		-- 设置过期时间
		if visitNum == 0 then
			redis.call('EXPIRE', key, expireTime)
		end
		-- 删除不在范围内的成员
		if visitNum > 1 then
			redis.call('ZREMRANGEBYSCORE', key, 1, beginScore)
		end

		return 1;
    `
	return script
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
