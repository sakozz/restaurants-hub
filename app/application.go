package app

import (
	"github.com/gin-gonic/gin"
	"resturants-hub.com/m/v2/database"
)

var (
	router = gin.Default()
)

func StartApplication() {
	database.RunMigrations()
	mapRoutes()

	router.Run(":3000")
}
