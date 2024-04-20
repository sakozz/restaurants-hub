package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"resturants-hub.com/m/v2/handlers"
	"resturants-hub.com/m/v2/middleware"
)

var (
	ssoHandler   handlers.SsoHandler   = handlers.NewSsoHandler()
	usersHandler handlers.UsersHandler = handlers.NewUsersHandler()
)

func mapRoutes() {
	router.GET("/", func(c *gin.Context) {
		router.LoadHTMLGlob("templates/*.html")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "This is an index page...",
		})
	})
	router.GET("/api/users", middleware.RequireAuth, usersHandler.List)
	router.GET("/api/users/:id", middleware.RequireAuth, usersHandler.Get)
	router.GET("/api/auth/:provider", ssoHandler.SsoLogin)
	router.GET("/api/auth/:provider/callback", ssoHandler.Callback)
	router.PUT("/api/auth/renew-session", middleware.RequireAuth, ssoHandler.RenewSession)
	router.POST("/api/auth/logout", middleware.RequireAuth, ssoHandler.Logout)
}
