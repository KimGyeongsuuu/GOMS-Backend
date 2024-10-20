package main

import (
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/router"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	fmt.Println("Connected to MySQL using GORM!")

	err = db.AutoMigrate(&model.Account{})
	if err != nil {
		log.Fatal("Failed to migrate table:", err)
	}

	if err != nil {
		log.Fatal("Failed to migrate table:", err)
	}

	r := router.SetupRouter(db)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}

}
