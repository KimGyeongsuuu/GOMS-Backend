package main

import (
	"GOMS-BACKEND-GO/controller"
	"GOMS-BACKEND-GO/database/cache"
	"GOMS-BACKEND-GO/database/mysql"
	"GOMS-BACKEND-GO/global/auth/jwt"
	"GOMS-BACKEND-GO/global/auth/jwt/middleware"
	"GOMS-BACKEND-GO/global/config"
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/repository"
	"GOMS-BACKEND-GO/service"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
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

	db, err := mysql.NewMySQLConnection()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	rdb = cache.NewRedisClient(ctx)

	err = db.AutoMigrate(&model.Account{}, &model.Outing{}, &model.Late{})
	if err != nil {
		log.Fatal("Failed to migrate tables:", err)
	}

	jwtConfig := config.JWT()
	outingConfig := config.Outing()

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
				"message": "goms server is running !!",
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
