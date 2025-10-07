package database

import (
	"News/service/internal/config"
	"context"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

const (
	mongoLocal    = "RichieMongo"
	redisLocal    = "RichieRedis"
	mariaLocal    = "RichieMaria"
	postgresLocal = "RichiePostgres"
)

type NewsOut struct {
	dig.Out
	MongoLocal    *mongo.Database `name:"mongo_news"`
	PostgresLocal *gorm.DB        `name:"postgres_news"`
	RedisLocal    *redis.Client   `name:"redis_news"`
	MariaLocal    *gorm.DB        `name:"maria_news"`
}

func NewNews(ctx context.Context, dbms config.DatabaseManageSystem) NewsOut {
	return NewsOut{
		MongoLocal:    newMongoDB(ctx, mongoLocal, dbms.MongoDBSystem[mongoLocal]),
		PostgresLocal: newPostgres(postgresLocal, dbms.PostgresSystem[postgresLocal]),
		RedisLocal:    newRedis(ctx, redisLocal, dbms.RedisSystem[redisLocal]),
		MariaLocal:    newMariaDB(mariaLocal, dbms.MariaDBSystem[mariaLocal]),
	}
}
