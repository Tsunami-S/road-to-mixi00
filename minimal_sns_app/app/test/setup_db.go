package test

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"log"
	"time"
	"minimal_sns_app/db"
	"testing"
)

func InitTestDB() *gorm.DB {
	dsn := os.Getenv("MYSQL_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ DB接続失敗: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}

func setupTestDB(t *testing.T) {
	db.DB = InitTestDB()
}

func setupTestDB_FOF(t *testing.T) {
	db.DB = InitTestDB()
}
