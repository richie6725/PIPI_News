package database

import (
	"News/service/internal/config"
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

func newRedis(ctx context.Context, dbName string, db config.Redis) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     (db.Host + ":" + db.Port),
		Password: db.Password,
		DB:       db.Database,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to ping Redis: %s, err: %v", dbName, err)
	}
	log.Printf("Connected to Redis: %s", dbName)

	return client
}
