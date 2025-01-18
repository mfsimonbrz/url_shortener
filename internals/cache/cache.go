package cache

import (
	"context"
	"time"
	"url_shortener/internals/models"

	"github.com/redis/go-redis/v9"
)

const REDIS_TTL = 45 //days

type CacheHandler struct {
	context context.Context
	client  *redis.Client
}

func NewCacheHandler(context context.Context, client *redis.Client) *CacheHandler {
	return &CacheHandler{context: context, client: client}
}

func (h *CacheHandler) Add(entry *models.Entry) error {
	ttl, err := time.ParseDuration("24h")
	if err != nil {
		return err
	}
	return h.client.Set(h.context, entry.ShortUrl, entry.Url, ttl*REDIS_TTL).Err()
}

func (h *CacheHandler) Get(key string) (string, error) {
	result, err := h.client.Get(h.context, key).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}
