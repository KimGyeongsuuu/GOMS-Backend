package main

import (
	"GOMS-BACKEND-GO/controller"
	"GOMS-BACKEND-GO/global/auth/jwt"
	"GOMS-BACKEND-GO/global/auth/jwt/middleware"
	"GOMS-BACKEND-GO/global/config"
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/repository"
	"GOMS-BACKEND-GO/service"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtProperties, jwtExpTimeProperties, err := config.LoadJwtProperties()
	if err != nil {
		log.Fatal("Failed to load JWT properties:", err)
	}

	outingProperties, err := config.LoadOutingProperties()
	if err != nil {
		log.Fatal("Failed to load Outing properties:", err)
	}

	db, err := setupDatabase()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	err = db.AutoMigrate(&model.Account{}, &model.Outing{})
	if err != nil {
		log.Fatal("Failed to migrate tables:", err)
	}

	rdb = setupRedis()

	refreshRepo := repository.NewRefreshTokenRepository(rdb)
	tokenAdapter := jwt.NewGenerateTokenAdapter(jwtProperties, jwtExpTimeProperties, rdb, refreshRepo)

	r := gin.Default()

	accountRepo := repository.NewAccountRepository(db)
	authUseCase := service.NewAuthService(accountRepo, tokenAdapter)
	authController := controller.NewAuthController(authUseCase)

	outingUUIDRepo := repository.NewOutingUUIDRepository(rdb, outingProperties)
	studentCouncilUseCase := service.NewStudentCouncilService(outingUUIDRepo)
	studentCouncilController := controller.NewStudentCouncilController(studentCouncilUseCase)

	outingRepo := repository.NewOutingRepository(db)
	outingUseCase := service.NewOutingService(outingRepo, accountRepo, outingUUIDRepo)
	outingController := controller.NewOutingController(outingUseCase)

	r.Use(middleware.AccountMiddleware(accountRepo, jwtProperties.AccessSecret))

	health := r.Group("/health")
	{
		health.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "goms server is running",
			})
		})
	}
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("signup", authController.SignUp)
		auth.POST("signin", authController.SignIn)

	}
	studentCouncil := r.Group("/api/v1/student-council")
	{
		studentCouncil.Use(middleware.AuthorizeRoleJWT(jwtProperties.AccessSecret, "ROLE_STUDENT_COUNCIL"))
		studentCouncil.POST("outing", studentCouncilController.CreateOuting)
	}
	outing := r.Group("/api/v1/outing")
	{
		outing.POST("/:outingUUID", outingController.OutingStudent)
		outing.GET("", outingController.ListOutingStudent)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func setupDatabase() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", user, password, host, port, database)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
