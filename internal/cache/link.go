package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IvanTime-Kai/url-shortener/internal/domain"
	"github.com/redis/go-redis/v9"
)

const linkTTL = 24 *time.Hour

type LinkCache struct {
	client *redis.Client
}

func NewLinkCache(client *redis.Client) *LinkCache {
	return &LinkCache{
		client: client,
	}
}

func (c *LinkCache) Get(ctx context.Context, code string) (*domain.Link, error) {
	val, err := c.client.Get(ctx, cacheKey(code)).Result()

	if err != nil {
		return  nil, err
	}

	var link domain.Link
	if err := json.Unmarshal([]byte(val), &link); err != nil {
		return nil, err
	}

	return &link, nil
}

func (c *LinkCache) Set(ctx context.Context, link *domain.Link) error {

	data, err := json.Marshal(link)

	if err != nil {
		return err
	}

	return c.client.Set(ctx, cacheKey(link.Code), data, linkTTL).Err()
}

func (c *LinkCache) Delete(ctx context.Context, code string) error {
	return  c.client.Del(ctx, cacheKey(code)).Err()
}

func cacheKey(code string) string {
    return "link:" + code
}