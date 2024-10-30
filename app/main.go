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

	outingBlackListProperties, err := config.LoadOutingBlackListProperties()
	if err != nil {
		log.Fatal("Failed to load Outing black list properties:", err)
	}

	db, err := setupDatabase()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	err = db.AutoMigrate(&model.Account{}, &model.Outing{}, &model.Late{})
	if err != nil {
		log.Fatal("Failed to migrate tables:", err)
	}

	rdb = setupRedis()

	refreshRepo := repository.NewRefreshTokenRepository(rdb)
	tokenAdapter := jwt.NewGenerateTokenAdapter(jwtProperties, jwtExpTimeProperties, rdb, refreshRepo)
	tokenParser := jwt.NewTokenParser()

	r := gin.Default()

	accountRepo := repository.NewAccountRepository(db)
	blackListRepo := repository.NewBlackListRepository(rdb)
	outingUUIDRepo := repository.NewOutingUUIDRepository(rdb, outingProperties)
	outingRepo := repository.NewOutingRepository(db)
	lateRepo := repository.NewLateRepository(db)
	authenticationRepo := repository.NewAuthenticationRepository(rdb)
	authCodeRepo := repository.NewAuthCodeRepository(rdb)

	authUseCase := service.NewAuthService(accountRepo, tokenAdapter, refreshRepo, tokenParser, authenticationRepo, authCodeRepo)
	outingUseCase := service.NewOutingService(outingRepo, accountRepo, outingUUIDRepo)
	lateUseCase := service.NewLateService(lateRepo)
	studentCouncilUseCase := service.NewStudentCouncilService(outingUUIDRepo, accountRepo, blackListRepo, outingBlackListProperties, outingRepo, lateRepo)
	accountUseCase := service.NewAccountService(accountRepo)

	authController := controller.NewAuthController(authUseCase)
	outingController := controller.NewOutingController(outingUseCase)
	lateController := controller.NewLateController(lateUseCase)
	studentCouncilController := controller.NewStudentCouncilController(studentCouncilUseCase)
	accountController := controller.NewAccountController(accountUseCase)

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
		auth.PATCH("", authController.TokenReissue)
		auth.POST("send/email", authController.SendAuthEmail)
		auth.GET("verify/email", authController.VerifyAuthCode)
	}
	studentCouncil := r.Group("/api/v1/student-council")
	{
		studentCouncil.Use(middleware.AuthorizeRoleJWT(jwtProperties.AccessSecret, "ROLE_STUDENT_COUNCIL"))
		studentCouncil.POST("outing", studentCouncilController.CreateOuting)
		studentCouncil.GET("accounts", studentCouncilController.FindOutingList)
		studentCouncil.GET("search", studentCouncilController.SearchAccountByInfo)
		studentCouncil.PATCH("authority", studentCouncilController.UpdateAuthority)
		studentCouncil.POST("black-list/:accountID", studentCouncilController.AddBlackList)
		studentCouncil.DELETE("black-list/:accountID", studentCouncilController.DeleteBlackList)
		studentCouncil.DELETE("outing/:accountID", studentCouncilController.DeleteOutingStudent)
		studentCouncil.GET("late", studentCouncilController.FindLateList)
	}
	outing := r.Group("/api/v1/outing")
	{
		outing.POST(":outingUUID", outingController.OutingStudent)
		outing.GET("", outingController.ListOutingStudent)
		outing.GET("count", outingController.CountOutingStudent)
		outing.GET("search", outingController.SearchOutingStudent)
	}
	late := r.Group("/api/v1/late")
	{
		late.GET("rank", lateController.GetLateStudentTop3)
	}
	account := r.Group("/api/v1/account")
	{
		account.DELETE("", accountController.WithDraw)
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
		Addr:     "svc.sel4.cloudtype.app:30258",
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
