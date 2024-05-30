package app

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"resturants-hub.com/m/v2/domains/restaurants"
	"resturants-hub.com/m/v2/domains/sessions"
	"resturants-hub.com/m/v2/domains/users"
	"resturants-hub.com/m/v2/middleware"
)

var (
	ssoHandler              sessions.SsoHandler                 = sessions.NewSsoHandler()
	usersHandler            users.UsersHandler                  = users.NewUsersHandler()
	adminRestaurantsHandler restaurants.AdminRestaurantsHandler = restaurants.NewAdminRestaurantsHandler()
)

func mapRoutes() {
	router.GET("/", func(c *gin.Context) {
		router.LoadHTMLGlob("templates/*.html")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "This is an index page...",
		})
	})

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"}
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "user-agent", "X-Requested-With", "Set-Cookie", "Cookie", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Accept-Language", "Accept-Encoding", "Accept", "Connection", "Host", "Referer", "Origin", "User-Agent"}
	router.Use(cors.New(config))

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
