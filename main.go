package main

import (
	"GOMS-BACKEND-GO/global/config"
	"GOMS-BACKEND-GO/global/jwt"
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/repository"
	"GOMS-BACKEND-GO/router"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var rdb *redis.Client
var ctx = context.Background()

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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", user, password, host, port, database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	fmt.Println("Connected to MySQL using GORM!")

	err = db.AutoMigrate(&model.Account{})
	err = db.AutoMigrate(&model.RefreshToken{})
	if err != nil {
		log.Fatal("Failed to migrate table:", err)
	}

	jwtProperties := &config.JwtProperties{
		AccessSecret:  []byte(os.Getenv("JWT_ACCESS_SECRET")),
		RefreshSecret: []byte(os.Getenv("JWT_REFRESH_SECRET")),
	}

	accessExp, err := strconv.Atoi(os.Getenv("JWT_ACCESS_EXP"))
	if err != nil {
		log.Fatal("Invalid access expiration time:", err)
	}
	refreshExp, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXP"))
	if err != nil {
		log.Fatal("Invalid refresh expiration time:", err)
	}

	jwtExpTimeProperties := &config.JwtExpTimeProperties{
		AccessExp:  int64(accessExp),
		RefreshExp: int64(refreshExp),
	}
	rdb = setupRedis()
	refreshRepo := repository.NewRefreshTokenRepository(db, rdb)

	tokenAdapter := jwt.NewGenerateTokenAdapter(jwtProperties, jwtExpTimeProperties, rdb, refreshRepo)

	r := router.SetupRouter(db, tokenAdapter)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}

}

func setupRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
		return nil
	}

	return rdb
}
