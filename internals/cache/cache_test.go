package cache

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"url_shortener/internals/models"

	"github.com/redis/go-redis/v9"
)

func setup() (*redis.Client, context.Context) {
	redisHost := os.Getenv("redisHost")
	redisPort := os.Getenv("redisPort")
	redisDb := os.Getenv("redisDb")

	ctx := context.Background()
	redisDatabase, _ := strconv.Atoi(redisDb)
	redis_client := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%s", redisHost, redisPort), DB: redisDatabase})

	return redis_client, ctx
}

func shutdown(context context.Context, redisClient *redis.Client) {
	redisClient.Del(context, "abc1234")
	defer redisClient.Conn().Close()
}

func TestAdd(t *testing.T) {
	redisClient, context := setup()
	defer shutdown(context, redisClient)
	cacheHandler := NewCacheHandler(context, redisClient)
	entry := &models.Entry{Url: "https://g1.globo.com/", ShortUrl: "abc1234"}
	err := cacheHandler.Add(entry)
	if err != nil {
		t.Error(err)
	}
	expected := "https://g1.globo.com/"
	redisCmd := redisClient.Get(context, "abc1234")
	got, err := redisCmd.Result()
	if err != nil {
		t.Error(err)
	}

	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestGet(t *testing.T) {
	redisClient, context := setup()
	defer shutdown(context, redisClient)
	redisClient.Set(context, "abc1234", "https://g1.globo.com/", 0)
	cacheHandler := NewCacheHandler(context, redisClient)
	expected := "https://g1.globo.com/"
	got, err := cacheHandler.Get("abc1234")
	if err != nil {
		t.Error(err)
	}

	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
