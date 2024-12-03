package main

import (
	"GOMS-BACKEND-GO/controller"
	"GOMS-BACKEND-GO/database/cache"
	"GOMS-BACKEND-GO/database/mongo"
	_ "GOMS-BACKEND-GO/docs"
	"GOMS-BACKEND-GO/global/auth/jwt"
	"GOMS-BACKEND-GO/global/auth/jwt/middleware"
	"GOMS-BACKEND-GO/global/config"
	"GOMS-BACKEND-GO/global/filter"
	"GOMS-BACKEND-GO/global/util"
	"GOMS-BACKEND-GO/repository"
	"GOMS-BACKEND-GO/service"
	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// db, err := mysql.NewMySQLConnection()
	// if err != nil {
	// 	log.Fatal("Failed to connect to the database:", err)
	// }

	rdb = cache.NewRedisClient(ctx)

	_, mongodb, err := mongo.NewMongoConnection()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// collectionNames := []string{"outings", "accounts", "lates"}
	// mongo.CreateCollections(mongodb, collectionNames)

	// err = db.AutoMigrate(&model.Account{}, &model.Outing{}, &model.Late{})
	// if err != nil {
	// 	log.Fatal("Failed to migrate tables:", err)
	// }

	errorFilter := filter.NewErrorFilter()

	jwtConfig := config.JWT()
	outingConfig := config.Outing()

	refreshRepo := repository.NewRefreshTokenRepository(rdb)
	token := jwt.NewToken(&jwtConfig, rdb, refreshRepo)
	passwordUtil := util.NewPasswordUtil()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        3600,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// mongo
	mongoAccountRepo := repository.NewMongoAccountRepository(mongodb)
	mongoOutingRepo := repository.NewMongoOutingRepository(mongodb)
	mongoLateRepo := repository.NewMongoLateRepository(mongodb)

	// mysql
	// accountRepo := repository.NewAccountRepository(db)
	// outingRepo := repository.NewOutingRepository(db)
	// lateRepo := repository.NewLateRepository(db)

	// redis
	blackListRepo := repository.NewBlackListRepository(rdb)
	outingUUIDRepo := repository.NewOutingUUIDRepository(rdb, &outingConfig)
	authenticationRepo := repository.NewAuthenticationRepository(rdb)
	authCodeRepo := repository.NewAuthCodeRepository(rdb)

	authUseCase := service.NewAuthService(mongoAccountRepo, token, token, refreshRepo, authenticationRepo, authCodeRepo, passwordUtil)
	outingUseCase := service.NewOutingService(mongoOutingRepo, mongoAccountRepo, outingUUIDRepo)
	lateUseCase := service.NewLateService(mongoLateRepo, mongoAccountRepo)
	studentCouncilUseCase := service.NewStudentCouncilService(outingUUIDRepo, mongoAccountRepo, blackListRepo, &outingConfig, mongoOutingRepo, mongoLateRepo)
	accountUseCase := service.NewAccountService(mongoAccountRepo)

	authController := controller.NewAuthController(authUseCase)
	outingController := controller.NewOutingController(outingUseCase)
	lateController := controller.NewLateController(lateUseCase)
	studentCouncilController := controller.NewStudentCouncilController(studentCouncilUseCase)
	accountController := controller.NewAccountController(accountUseCase)

	r.Use(middleware.AccountMiddleware(mongoAccountRepo, []byte(jwtConfig.AccessSecret)))
	r.Use(errorFilter.Register())

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
		studentCouncil.GET("accounts", studentCouncilController.FindAccountList)
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
