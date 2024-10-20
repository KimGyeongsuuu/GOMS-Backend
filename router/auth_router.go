package router

import (
	"GOMS-BACKEND-GO/controller"
	"GOMS-BACKEND-GO/global/jwt"
	"GOMS-BACKEND-GO/repository"
	"GOMS-BACKEND-GO/service"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, tokenAdapter *jwt.GenerateTokenAdapter) *gin.Engine {

	r := gin.Default()

	accountRepo := repository.NewAccountRepository(db)
	authUseCase := service.NewAuthService(accountRepo, tokenAdapter)
	authController := controller.NewAuthController(authUseCase)

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

	return r
}
