package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

func InitDB() {
	user := "postgres"
	pass := ""
	host := "localhost"
	port := "5432"
	dbname := "postgres"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, pass, dbname, port)

	var err interface{}
	GormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Database connected successfully!")

	sqlDB, err := GormDB.DB()
	if err != nil {
		fmt.Println("Error fetching database object:", err)
		return
	}

	sqlDB.SetMaxIdleConns(3)
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
}

func CloseDB() {
	if GormDB == nil {
		fmt.Println("No active database connection to close.")
		return
	}

	sqlDB, err := GormDB.DB()
	if err != nil {
		fmt.Println("Error while fetching database object:", err)
		return
	}

	err = sqlDB.Close()
	if err != nil {
		fmt.Println("Error while closing database connection:", err)
	} else {
		fmt.Println("Database connection closed successfully!")
	}
}
