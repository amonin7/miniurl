package ratelimit

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const namespace = "rl"

func NewLimiter(client *redis.Client, action string, period time.Duration, limit int64) *Limiter {
	return &Limiter{
		client: client,
		action: action,
		period: period,
		limit:  limit,
	}
}

type Limiter struct {
	client *redis.Client

	action string
	period time.Duration
	limit  int64
}

//go:embed incr_inspirenx.lua
var IncrExpireLua string
var IncrExpireScript = redis.NewScript(IncrExpireLua)

func (l *Limiter) key(ts time.Time) string {
	interval := ts.UTC().UnixNano() / l.period.Nanoseconds()
	return fmt.Sprintf("%s:%s:%x", namespace, l.action, interval)
}

func (l *Limiter) CanDoAt(ctx context.Context, ts time.Time) (bool, error) {
	key := l.key(ts)
	ttlMs := l.period.Milliseconds()

	rawCount, err := IncrExpireScript.Run(ctx, l.client, []string{key}, ttlMs).Result()
	if err != nil {
		return false, err
	}
	count := rawCount.(int64)

	return count <= l.limit, nil
}
