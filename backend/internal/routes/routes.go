package routes

import (
	"employee-platform/backend/config"
	"employee-platform/backend/internal/controllers"
	"employee-platform/backend/internal/middleware"
	"employee-platform/backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(router *gin.Engine, db *mongo.Database, cfg *config.Config) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	mailerService := services.NewMailerService(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom)
	authService := services.NewAuthService(db, cfg.JWTSecret, mailerService)
	employeeService := services.NewEmployeeService(db, cfg.FrontendURL, mailerService)
	profileService := services.NewProfileService(db)

	authController := controllers.NewAuthController(authService)
	employeeController := controllers.NewEmployeeController(employeeService)
	profileController := controllers.NewProfileController(profileService)

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/admin-signup/send-otp", authController.SendAdminSignupOTP)
			auth.POST("/admin-signup/verify-otp", authController.VerifyAdminSignupOTP)
			auth.POST("/login/send-otp", authController.SendLoginOTP)
			auth.POST("/login/verify-otp", authController.VerifyLoginOTP)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			protected.GET("/auth/me", authController.Me)
			protected.PUT("/profile/me", profileController.CompleteProfile)

			adminOnly := protected.Group("")
			adminOnly.Use(middleware.RequireRoles("admin", "super_admin", "manager"))
			{
				adminOnly.POST("/employees", employeeController.CreateEmployee)
				adminOnly.GET("/employees", employeeController.ListEmployees)
			}
		}
	}
}