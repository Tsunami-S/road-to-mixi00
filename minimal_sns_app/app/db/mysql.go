package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"minimal_sns_app/model"
	"os"
	"time"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("❌ MYSQL_DSN is not set in environment")
	}

	var err error

	for i := 1; i <= 10; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err == nil {
			log.Println("✅ Successfully connected to the database")

			err = DB.AutoMigrate(
				&model.User{},
				&model.BlockList{},
			)
			if err != nil {
				log.Fatalf("❌ AutoMigrate failed: %v", err)
			}

			sqlDB, err := DB.DB()
			if err != nil {
				log.Fatalf("❌ Failed to retrieve underlying sql.DB instance: %v", err)
			}

			sqlDB.SetMaxOpenConns(10)
			sqlDB.SetMaxIdleConns(5)
			sqlDB.SetConnMaxLifetime(time.Hour)

			return
		}

		log.Printf("⚠️ Failed to connect to the database (attempt %d): %v", i, err)
		time.Sleep(3 * time.Second)
	}

	log.Fatalf("❌ Failed to connect to the database after 10 attempts: %v", err)
}
