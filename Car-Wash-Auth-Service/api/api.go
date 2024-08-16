package api

import (
	"github.com/exam-5/Car-Wash-Auth-Service/api/handler"
	middleware "github.com/exam-5/Car-Wash-Auth-Service/api/midlleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Auth Service
// @version 1.0
// @description Auth Service
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewGin(h handler.Handler) *gin.Engine {
	r := gin.Default()

	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(middleware.MiddleWare())
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/logout", h.Logout)
		auth.POST("/forgot", h.ForgotPassword)
		auth.POST("/reset", h.ResetPassword)
	}


	user := r.Group("/user")
	{
		user.GET("/profile/:id", h.GetProfile)
		user.PUT("/profile", h.UpdateProfile)
		user.PUT("/password", h.ChangePassword)
		user.GET("/all", h.GetAllUsers)
	}

	return r
}
// write swagger