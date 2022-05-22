package test

import (
	"fmt"
	rateLimit "rate-limit"
	"rate-limit/redis"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var limit = rateLimit.New(&redis.Config{Address: "127.0.0.1:6379"})

func TestDefaultLimiterTimes(t *testing.T) {
	createdDate, _ := strconv.Atoi(time.Now().Format("20060102150405"))

	var GetRandomKey = func(key string) string {
		return fmt.Sprintf("test:%s:%d", key, createdDate)
	}

	tests := []struct {
		name       string
		key        string
		IsLimitErr error
	}{
		{
			name: "测试 RateLimiter-1",
			key:  GetRandomKey("test-DefaultLimiter"),

			IsLimitErr: nil,
		},
		{
			name:       "测试 RateLimiter-2",
			key:        GetRandomKey("test-DefaultLimiter"),
			IsLimitErr: nil,
		},
		{
			name:       "测试 RateLimiter-3",
			key:        GetRandomKey("test-DefaultLimiter"),
			IsLimitErr: nil,
		},
		{
			name:       "测试 RateLimiter-4",
			key:        GetRandomKey("test-DefaultLimiter"),
			IsLimitErr: rateLimit.RateLimitErr,
		},
		{
			name:       "测试 RateLimiter-5",
			key:        GetRandomKey("test-DefaultLimiter"),
			IsLimitErr: rateLimit.RateLimitErr,
		},
		{
			name:       "测试 RateLimiter-6",
			key:        GetRandomKey("test-DefaultLimiter-1"),
			IsLimitErr: nil,
		},
	}
	for _, item := range tests {
		t.Run(item.name, func(t *testing.T) {
			got := limit.RateLimiter(
				rateLimit.ExpireTime(10),
				rateLimit.MaxThreads(3),
				rateLimit.IsLimitTime(true),
				rateLimit.Key(item.key),
			)
			assert.Equal(t, item.IsLimitErr, got)
		})
	}
}

func TestDefaultLimiterNotTimes(t *testing.T) {
	createdDate, _ := strconv.Atoi(time.Now().Format("20060102150405"))

	var GetRandomKey = func(key string) string {
		return fmt.Sprintf("test:%s:%d", key, createdDate)
	}

	tests := []struct {
		name       string
		key        string
		IsLimitErr error
	}{
		{
			name: "测试 RateLimiter-1",
			key:  GetRandomKey("test-DefaultLimiterNotTimes"),

			IsLimitErr: nil,
		},
		{
			name:       "测试 RateLimiter-2",
			key:        GetRandomKey("test-DefaultLimiterNotTimes"),
			IsLimitErr: nil,
		},
		{
			name:       "测试 RateLimiter-3",
			key:        GetRandomKey("test-DefaultLimiterNotTimes"),
			IsLimitErr: nil,
		},
		{
			name:       "测试 RateLimiter-4",
			key:        GetRandomKey("test-DefaultLimiterNotTimes"),
			IsLimitErr: nil,
		},
		{
			name:       "测试 RateLimiter-5",
			key:        GetRandomKey("test-DefaultLimiterNotTimes"),
			IsLimitErr: nil,
		},
		{
			name:       "测试 RateLimiter-6",
			key:        GetRandomKey("test-DefaultLimiterNotTimes"),
			IsLimitErr: nil,
		},
	}
	for _, item := range tests {
		t.Run(item.name, func(t *testing.T) {
			got := limit.RateLimiter(
				rateLimit.ExpireTime(10),
				rateLimit.MaxThreads(3),
				rateLimit.Key(item.key),
			)
			assert.Equal(t, item.IsLimitErr, got)
		})
	}
}

func TestTimesLimiter(t *testing.T) {
	createdDate, _ := strconv.Atoi(time.Now().Format("20060102150405"))

	var GetRandomKey = func(key string) string {
		return fmt.Sprintf("test:%s:%d", key, createdDate)
	}

	tests := []struct {
		name       string
		key        string
		IsLimitErr error
	}{
		{
			name: "测试 TimesLimiter-1",
			key:  GetRandomKey("test-TimesLimiter"),

			IsLimitErr: nil,
		},
		{
			name:       "测试 TimesLimiter-2",
			key:        GetRandomKey("test-TimesLimiter"),
			IsLimitErr: nil,
		},
		{
			name:       "测试 TimesLimiter-3",
			key:        GetRandomKey("test-TimesLimiter"),
			IsLimitErr: nil,
		},
		{
			name:       "测试 TimesLimiter-4",
			key:        GetRandomKey("test-TimesLimiter"),
			IsLimitErr: rateLimit.RateLimitErr,
		},
		{
			name:       "测试 TimesLimiter-5",
			key:        GetRandomKey("test-TimesLimiter"),
			IsLimitErr: rateLimit.RateLimitErr,
		},
		{
			name:       "测试 TimesLimiter-6",
			key:        GetRandomKey("test-TimesLimiter"),
			IsLimitErr: rateLimit.RateLimitErr,
		},
	}
	for _, item := range tests {
		t.Run(item.name, func(t *testing.T) {
			got := limit.TimesLimiter(item.key, 3, 10)
			assert.Equal(t, item.IsLimitErr, got)
		})
	}
}
