package rateLimit

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	DefaultExpireTime = 10 // 单位秒
	DefaultMaxThreads = 1
	DefaultPrefix     = "test"
)

type LimitClient struct {
	rateLimit redis.redis
}

func New(conf *redis.redis.Config) *LimitClient {
	return &LimitClient{
		rateLimit: redis.redis.New(conf),
	}
}

func (p *LimitClient) RateLimiter(param ...Param) gin.HandlerFunc {
	return func(c *gin.Context) {
		ps := evaluateParam(param)

		validAndAssignInput(c, ps)

		pipe := p.rateLimit.Pipeline(c)
		pipe.Send("INCR", ps.Key)
		pipe.Send("TTL", ps.Key)

		replies, err := pipe.Receive()
		if err != nil {
			log.Errorw("pipe filed", "message", ps.Key, "err", err)
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"ok":      false,
				"ecode":   ecode.RateLimitInvalid,
				"code":    "RATE_LIMIT",
				"message": "请稍后重试！",
			})
			return
		}

		var (
			current = replies[0].(int64)
			ttl     = replies[1].(int64)
		)

		if current == int64(1) || ttl == int64(-1) {
			p.rateLimit.Do(c, "EXPIRE", ps.Key, ps.ExpireTime)
		}

		if current > ps.MaxThreads {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"ok":      false,
				"ecode":   ecode.RateLimitInvalid,
				"code":    "RATE_LIMIT",
				"message": "请稍后重试！",
			})
			return

		}
		c.Next()

		if !ps.IsLimitTime {
			defer p.rateLimit.Del(c, ps.Key)
		}
	}
}

func validAndAssignInput(ctx *gin.Context, p *params) {
	keyItem := ""
	userID, exists := ctx.Get("userKey")
	if exists && userID != "" {
		keyItem = userID.(string)
		if p.IsLimitUser {
			p.Key += fmt.Sprintf("%s:limit:user:all:%s", DefaultPrefix, keyItem)
		}
	}

	if p.ExpireTime == 0 {
		p.ExpireTime = DefaultExpireTime
	}

	if p.MaxThreads == 0 {
		p.MaxThreads = DefaultMaxThreads
	}

	if p.Key == "" {
		// 格式 rlg:60:POST:gold:/gold/issueGold:118.112.12.34
		p.Key = fmt.Sprintf("%s:%s:%s:%s", DefaultPrefix, ctx.Request.Method, ctx.Request.URL.Path, keyItem)
	}
}
