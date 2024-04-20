package app

import (
	"github.com/gin-gonic/gin"
	"resturants-hub.com/m/v2/database"
	"resturants-hub.com/m/v2/domains/sso"
)

var (
	router = gin.Default()
)

func StartApplication() {
	database.RunMigrations()
	sso.Setup()
	mapRoutes()

	router.Run(":3000")
}
