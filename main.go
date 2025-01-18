package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"url_shortener/internals/cache"
	"url_shortener/internals/data"
	"url_shortener/internals/handler"
	"url_shortener/internals/web"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	dbHost := os.Getenv("dbHost")
	dbPort := os.Getenv("dbPort")
	dbUser := os.Getenv("dbUser")
	dbPass := os.Getenv("dbPass")
	dbName := os.Getenv("dbName")
	redisHost := os.Getenv("redisHost")
	redisPort := os.Getenv("redisPort")
	redisDb := os.Getenv("redisDb")
	appHost := os.Getenv("appHost")
	appPort := os.Getenv("appPort")

	ctx := context.Background()
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName))
	if err != nil {
		log.Fatal(err)
	}
	redisDatabase, _ := strconv.Atoi(redisDb)
	redis_client := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%s", redisHost, redisPort), DB: redisDatabase})

	entryData := data.NewEntryData(db)
	entryData.InitDB()
	cacheHandler := cache.NewCacheHandler(ctx, redis_client)
	entryHandler := handler.NewEntryHandler(entryData, cacheHandler)
	webHandler := web.NewEntryWebHandler(entryHandler)

	defer db.Close()
	defer redis_client.Conn().Close()

	r := gin.Default()
	r.GET("/", webHandler.HealthCheck)
	r.GET("/:short_url", webHandler.GetEntry)
	r.POST("/", webHandler.AddUrlEntry)

	r.Run(fmt.Sprintf("%s:%s", appHost, appPort))
}
