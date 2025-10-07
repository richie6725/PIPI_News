package database

import (
	"News/service/internal/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

func newPostgres(dbName string, db config.Postgres) *gorm.DB {
	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Taipei",
		db.Host, db.User, db.Password, db.Database, db.Port)

	gormDB, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}
	if db.MaxIdle == 0 || db.MaxOpen == 0 {
		log.Fatalf("%s missing maxIdle or maxOpen", dbName)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(db.MaxIdle)
	sqlDB.SetMaxOpenConns(db.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("Error pinging database %s: %v", dbName, err)
	}
	log.Printf("Pinged successfully postgres database: %s", dbName)

	return gormDB
}
