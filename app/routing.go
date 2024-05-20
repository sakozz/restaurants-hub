package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"resturants-hub.com/m/v2/handlers"
	"resturants-hub.com/m/v2/middleware"
)

var (
	ssoHandler              handlers.SsoHandler              = handlers.NewSsoHandler()
	usersHandler            handlers.UsersHandler            = handlers.NewUsersHandler()
	adminRestaurantsHandler handlers.AdminRestaurantsHandler = handlers.NewAdminRestaurantsHandler()
)

func mapRoutes() {
	router.GET("/", func(c *gin.Context) {
		router.LoadHTMLGlob("templates/*.html")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "This is an index page...",
		})
	})

	/* Admin routes */
	adminRoutes := router.Group("/api/admin", middleware.RequireAuth)
	{
		/* Admin Restaurant routes */
		adminRestaurantsRoutes := adminRoutes.Group("/restaurants")
		adminRestaurantsRoutes.POST("/", adminRestaurantsHandler.Create)
		adminRestaurantsRoutes.GET("/", adminRestaurantsHandler.List)
		adminRestaurantsRoutes.GET("/:id", adminRestaurantsHandler.Get)
		adminRestaurantsRoutes.PUT("/:id", adminRestaurantsHandler.Update)
		adminRestaurantsRoutes.PATCH("/:id", adminRestaurantsHandler.Update)

		adminUsersRoutes := adminRoutes.Group("/users")
		adminUsersRoutes.GET("/", usersHandler.List)
		adminUsersRoutes.GET("/profile", usersHandler.Profile)
		adminUsersRoutes.GET("/:id", usersHandler.Get)
	}

	/* Auth routes */
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.GET("/:provider", ssoHandler.SsoLogin)
		authRoutes.GET("/:provider/callback", ssoHandler.Callback)
		authRoutes.PUT("/renew-session", middleware.RequireAuth, ssoHandler.RenewSession)
		authRoutes.POST("/logout", middleware.RequireAuth, ssoHandler.Logout)
	}
}
