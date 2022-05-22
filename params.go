package rateLimit

import "errors"

const (
	DefaultExpireTime = 10 // 单位秒
	DefaultMaxThreads = 1
	DefaultPrefix     = "test"
)

var RateLimitErr = errors.New("请求过于频繁了,请稍后再试!")

type params struct {
	Key         string `json:"key"`
	MaxThreads  int64  `json:"maxThreads"`  // 最大的线程数
	ExpireTime  int64  `json:"expireTime"`  // 到期时间，秒
	IsLimitTime bool   `json:"isLimitTime"` // 是否在一定时间内限速
	IsLimitUser bool   `json:"isLimitUser"` // 是否限制用户，使用这个用户所有的接口只能访问对应的次数，并且一定要先验证用户是否登陆，慎用
}

type Param func(*params)

func evaluateParam(param []Param) *params {
	ps := &params{}

	for _, p := range param {
		p(ps)
	}
	return ps
}

func Key(key string) Param {
	return func(o *params) {
		o.Key = key
	}
}

func MaxThreads(maxThreads int64) Param {
	return func(o *params) {
		o.MaxThreads = maxThreads
	}
}

func ExpireTime(expireTime int64) Param {
	return func(o *params) {
		o.ExpireTime = expireTime
	}
}

func IsLimitTime(limit bool) Param {
	return func(o *params) {
		o.IsLimitTime = limit
	}
}

func IsLimitUser(limit bool) Param {
	return func(o *params) {
		o.IsLimitUser = limit
	}
}
