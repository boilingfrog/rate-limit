package test

import (
	"context"
	"fmt"
	rateLimit "rate-limit"
	"rate-limit/redis"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var limit = rateLimit.New(&redis.Config{Address: "127.0.0.1:6379"})

func TestRateLimiterTimes(t *testing.T) {
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
			got := limit.RateLimiter(context.Background(),
				rateLimit.ExpireTime(10),
				rateLimit.MaxThreads(3),
				rateLimit.IsLimitTime(true),
				rateLimit.Key(item.key),
			)
			assert.Equal(t, item.IsLimitErr, got)
		})
	}
}

func TestRateLimiterNotTimes(t *testing.T) {
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
			got := limit.RateLimiter(context.Background(),
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
			got := limit.TimesLimiter(context.Background(), item.key, 3, 10)
			assert.Equal(t, item.IsLimitErr, got)
		})
	}
}

func TestUserTimesLimiter(t *testing.T) {
	createdDate, _ := strconv.Atoi(time.Now().Format("20060102150405"))

	var GetRandomKey = func(key string) string {
		return fmt.Sprintf("test:%s:%d", key, createdDate)
	}

	tests := []struct {
		name       string
		User       string
		IsLimitErr error
	}{
		{
			name: "测试 UserTimes-1",
			User: GetRandomKey("test-UserTimes"),

			IsLimitErr: nil,
		},
		{
			name:       "测试 UserTimes-2",
			User:       GetRandomKey("test-UserTimes"),
			IsLimitErr: nil,
		},
		{
			name:       "测试 UserTimes-3",
			User:       GetRandomKey("test-UserTimes"),
			IsLimitErr: nil,
		},
		{
			name:       "测试 UserTimes-4",
			User:       GetRandomKey("test-UserTimes"),
			IsLimitErr: rateLimit.RateLimitErr,
		},
		{
			name:       "测试 UserTimes-5",
			User:       GetRandomKey("test-UserTimes"),
			IsLimitErr: rateLimit.RateLimitErr,
		},
		{
			name:       "测试 UserTimes-6",
			User:       GetRandomKey("test-UserTimes"),
			IsLimitErr: rateLimit.RateLimitErr,
		},
	}
	for _, item := range tests {
		t.Run(item.name, func(t *testing.T) {
			got := limit.UserTimesLimiter(context.Background(), item.User, 3, 10)
			assert.Equal(t, item.IsLimitErr, got)
		})
	}
}

func TestUserSingleRequestLimiter(t *testing.T) {
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
			name: "测试 SingleRequestLimiter-1",
			key:  GetRandomKey("test-SingleRequestLimiter"),

			IsLimitErr: nil,
		},
		{
			name:       "测试 SingleRequestLimiter-2",
			key:        GetRandomKey("test-SingleRequestLimiter"),
			IsLimitErr: nil,
		},
		{
			name:       "测试 SingleRequestLimiter-3",
			key:        GetRandomKey("test-SingleRequestLimiter"),
			IsLimitErr: nil,
		},
		{
			name:       "测试 SingleRequestLimiter-4",
			key:        GetRandomKey("test-SingleRequestLimiter"),
			IsLimitErr: nil,
		},
	}
	for _, item := range tests {
		t.Run(item.name, func(t *testing.T) {
			got := limit.UserSingleRequestLimiter(context.Background(), item.key, 10)
			assert.Equal(t, item.IsLimitErr, got)
		})
	}
}
