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

	authService := services.NewAuthService(db, cfg.JWTSecret)
	employeeService := services.NewEmployeeService(db)
	profileService := services.NewProfileService(db)

	authController := controllers.NewAuthController(authService)
	employeeController := controllers.NewEmployeeController(employeeService)
	profileController := controllers.NewProfileController(profileService)

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/admin-signup", authController.AdminSignup)
			auth.POST("/login", authController.Login)
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