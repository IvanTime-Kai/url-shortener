package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	ratelimit = 10
	ratePeriod = 1 * time.Minute
)

type RateLimit struct {
	client *redis.Client
}

func NewRateLimit(client *redis.Client) *RateLimit {
	return &RateLimit{
		client: client,
	}
}

func (r *RateLimit) Allow(ctx context.Context, ip string) (bool, error) {
	key := fmt.Sprintf("rate:%s", ip)

	// Increment counter for ip
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		r.client.Expire(ctx, key, ratePeriod)
	}

	return count <= ratelimit, nil
}