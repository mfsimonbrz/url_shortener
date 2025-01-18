package handler

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"testing"
	"url_shortener/internals/cache"
	"url_shortener/internals/data"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func setup() (*sql.DB, *redis.Client, context.Context, error) {
	dbHost := os.Getenv("testDbHost")
	dbPort := os.Getenv("testDbPort")
	dbUser := os.Getenv("testDbUser")
	dbPass := os.Getenv("testDbPass")
	dbName := os.Getenv("testDbName")
	redisHost := os.Getenv("redisHost")
	redisPort := os.Getenv("redisPort")
	redisDb := os.Getenv("redisDb")

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName))
	if err != nil {
		return nil, nil, nil, err
	}

	ctx := context.Background()
	redisDatabase, _ := strconv.Atoi(redisDb)
	redis_client := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%s", redisHost, redisPort), DB: redisDatabase})

	return db, redis_client, ctx, nil
}

func shutdown(db *sql.DB, redisClient *redis.Client, context context.Context) error {
	_, err := db.Exec("DELETE FROM public.entries")
	if err != nil {
		return err
	}

	defer db.Close()

	redisClient.Del(context, "abc1234")
	defer redisClient.Conn().Close()

	return nil
}

func TestAddUrlEntry(t *testing.T) {
	db, redisClient, context, err := setup()
	defer shutdown(db, redisClient, context)

	if err != nil {
		t.Error(err)
	}

	entryData := data.NewEntryData(db)
	cacheHandler := cache.NewCacheHandler(context, redisClient)
	entryHandler := NewEntryHandler(entryData, cacheHandler)
	newEntry, err := entryHandler.AddUrlEntry("https://g1.globo.com/")

	if err != nil {
		t.Error(err)
	}

	if newEntry.Url == "" {
		t.Error("should have url")
	}

	newEntry, err = entryHandler.AddUrlEntry("https://g1.globo.com/")

	if err != nil {
		t.Error(err)
	}

	if newEntry.ID == 0 {
		t.Error("should be greater than zero")
	}
}

func TestRetrieveUrl(t *testing.T) {
	db, redisClient, context, err := setup()
	defer shutdown(db, redisClient, context)

	if err != nil {
		t.Error(err)
	}

	entryData := data.NewEntryData(db)
	cacheHandler := cache.NewCacheHandler(context, redisClient)
	entryHandler := NewEntryHandler(entryData, cacheHandler)

	newEntry, err := entryHandler.AddUrlEntry("https://g1.globo.com/")

	if err != nil {
		t.Error(err)
	}

	got, err := entryHandler.RetrieveUrl(newEntry.ShortUrl)
	expected := "https://g1.globo.com/"

	if err != nil {
		t.Error(err)
	}

	if got != expected {
		t.Errorf("got %q, expeted %q", got, expected)
	}
}
