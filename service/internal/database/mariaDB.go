package database

import (
	"News/service/internal/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

func newMariaDB(dbName string, db config.MariaDB) *gorm.DB {
	dataSourceName :=
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&readTimeout=30s&writeTimeout=30s",
			db.Account, db.Password, db.Host, db.Port, db.Database)

	gormDB, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{
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
	log.Printf("Pinged successfully maria database: %s", dbName)

	return gormDB
}
