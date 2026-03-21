package main

import (
	"log"

	"employee-platform/backend/config"
	"employee-platform/backend/internal/database"
	"employee-platform/backend/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	mongoConn, err := database.NewMongo(cfg)
	if err != nil {
		log.Fatal("failed to connect to mongodb:", err)
	}

	r := gin.Default()
	routes.Setup(r, mongoConn.DB, cfg)

	log.Println("server running on :" + cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}