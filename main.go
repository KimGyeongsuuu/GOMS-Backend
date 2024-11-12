package main

import (
	"GOMS-BACKEND-GO/controller"
	"GOMS-BACKEND-GO/global/auth/jwt"
	"GOMS-BACKEND-GO/global/auth/jwt/middleware"
	"GOMS-BACKEND-GO/global/config"
	"GOMS-BACKEND-GO/repository"
	"GOMS-BACKEND-GO/service"
	"context"
	"fmt"
	"log"
	"time"

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
		log.Fatal("Error loading .env file ", err.Error())
	}

	if err := config.Load("./resource/app.yml"); err != nil {
		log.Fatal(err.Error())
	}

	jwtConfig := config.JWT()
	outingConfig := config.Outing()

	db, err := setupDatabase()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	rdb = setupRedis()

	refreshRepo := repository.NewRefreshTokenRepository(rdb)
	tokenAdapter := jwt.NewGenerateTokenAdapter(&jwtConfig, rdb, refreshRepo)
	tokenParser := jwt.NewTokenParser()

	r := gin.Default()

	accountRepo := repository.NewAccountRepository(db)
	blackListRepo := repository.NewBlackListRepository(rdb)
	outingUUIDRepo := repository.NewOutingUUIDRepository(rdb, &outingConfig)
	outingRepo := repository.NewOutingRepository(db)
	lateRepo := repository.NewLateRepository(db)
	authenticationRepo := repository.NewAuthenticationRepository(rdb)
	authCodeRepo := repository.NewAuthCodeRepository(rdb)

	authUseCase := service.NewAuthService(accountRepo, tokenAdapter, refreshRepo, tokenParser, authenticationRepo, authCodeRepo)
	outingUseCase := service.NewOutingService(outingRepo, accountRepo, outingUUIDRepo)
	lateUseCase := service.NewLateService(lateRepo)
	studentCouncilUseCase := service.NewStudentCouncilService(outingUUIDRepo, accountRepo, blackListRepo, &outingConfig, outingRepo, lateRepo)
	accountUseCase := service.NewAccountService(accountRepo)

	authController := controller.NewAuthController(authUseCase)
	outingController := controller.NewOutingController(outingUseCase)
	lateController := controller.NewLateController(lateUseCase)
	studentCouncilController := controller.NewStudentCouncilController(studentCouncilUseCase)
	accountController := controller.NewAccountController(accountUseCase)

	r.Use(middleware.AccountMiddleware(accountRepo, []byte(jwtConfig.AccessSecret)))

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
		studentCouncil.Use(middleware.AuthorizeRoleJWT([]byte(jwtConfig.AccessSecret), "ROLE_STUDENT_COUNCIL"))
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
	user := config.Data().Mysql.User
	password := config.Data().Mysql.Pass
	host := config.Data().Mysql.Host
	port := config.Data().Mysql.Port
	database := config.Data().Mysql.Db

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true", user, password, host, port, database)

	var db *gorm.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("failed to connect to the database: %w", err)
}

func setupRedis() *redis.Client {
	host := config.Data().Redis.Host
	port := config.Data().Redis.Port

	addr := fmt.Sprintf("%s:%d", host, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
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
