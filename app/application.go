package app

import (
	"os"

	"github.com/gin-gonic/gin"
	"resturants-hub.com/m/v2/database"
)

var (
	router = gin.Default()
)

func StartApplication() {
	database.RunMigrations()
	mapRoutes()

	// Use env variable for port configuration if available, default to 3000 otherwise
	if port := os.Getenv("PORT"); port != "" {
		router.Run(":" + port)
	} else {
		router.Run(":3000")
	}
}
