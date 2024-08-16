package api

import (
	"github.com/exam-5/Car-Wash-Api-Gateway/api/handler"
	"github.com/exam-5/Car-Wash-Api-Gateway/api/middleware"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Api Gateway
// @version 1.0
// @description Api Gateway Booking Service
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewGin() *gin.Engine {
	h := handler.NewHandler()
	r := gin.Default()
	r.GET("api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(middleware.MiddleWare())
	r.Use(middleware.CasbinMiddleware(h.Client.Enforcer))

	r.POST("/services", h.CreateService)
	r.GET("/services", h.GetServices)
	r.GET("/services/:id", h.GetService)
	r.PUT("/services/:id", h.UpdateService)
	r.DELETE("/services/:id", h.DeleteService)
	r.GET("/services/search", h.SearchServices)
	r.GET("/services/popular", h.PopularServices)

	r.POST("/providers", h.CreateProvider)
	r.GET("/providers", h.GetProviders)
	r.GET("/providers/:id", h.GetProvider)
	r.PUT("/providers/:id", h.UpdateProvider)
	r.DELETE("/providers/:id", h.DeleteProvider)

	r.GET("/providers/search", h.SearchProviders)

	r.POST("/bookings", h.CreateBooking)
	r.GET("/bookings", h.GetBookings)
	r.GET("/bookings/:id", h.GetBooking)
	r.PUT("/bookings/:id", h.UpdateBooking)
	r.DELETE("/bookings/:id", h.DeleteBooking)
	r.PUT("/bookings/:id/confirm", h.ConfirmBooking)

	r.POST("/reviews", h.CreateReview)
	r.GET("/reviews", h.GetReviews)
	r.GET("/reviews/:id", h.GetReview)
	r.PUT("/reviews/:id", h.UpdateReview)
	r.DELETE("/reviews/:id", h.DeleteReview)

	r.POST("/payments", h.CreatePayment)
	r.GET("/payments", h.ListPayments)
	r.GET("/payments/:id", h.GetPayment)

	r.GET("/notifications/:booking_id", h.GetNotigication)

	return r
}
