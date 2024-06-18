package app

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"resturants-hub.com/m/v2/domains/invitations"
	"resturants-hub.com/m/v2/domains/restaurants"
	"resturants-hub.com/m/v2/domains/sessions"
	"resturants-hub.com/m/v2/domains/users"
	"resturants-hub.com/m/v2/middleware"
)

var (
	ssoHandler         sessions.SsoHandler            = sessions.NewSsoHandler()
	usersHandler       users.UsersHandler             = users.NewUsersHandler()
	invitationsHandler invitations.InvitationsHandler = invitations.NewInvitationsHandler()
	restaurantsHandler restaurants.RestaurantsHandler = restaurants.NewAdminRestaurantsHandler()
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
		adminRestaurantsRoutes.POST("/", restaurantsHandler.Create)
		adminRestaurantsRoutes.GET("/", restaurantsHandler.List)
		adminRestaurantsRoutes.GET("/:id", restaurantsHandler.Get)
		adminRestaurantsRoutes.PUT("/:id", restaurantsHandler.Update)
		adminRestaurantsRoutes.PATCH("/:id", restaurantsHandler.Update)

		adminUsersRoutes := adminRoutes.Group("/users")
		adminUsersRoutes.POST("/", usersHandler.Create)
		adminUsersRoutes.GET("/", usersHandler.List)
		adminUsersRoutes.GET("/profile", usersHandler.Profile)
		adminUsersRoutes.GET("/:id", usersHandler.Get)

		/* Admin Invitations routes */
		adminInvitationsRoutes := adminRoutes.Group("/invitations")
		adminInvitationsRoutes.POST("/", invitationsHandler.Create)
		adminInvitationsRoutes.GET("/", invitationsHandler.List)
		adminInvitationsRoutes.GET("/:id", invitationsHandler.Get)
		adminInvitationsRoutes.PATCH("/:id", invitationsHandler.Update)
	}

	/* Manager's Restaurant routes */
	restaurantsRoutes := router.Group("/api/my-restaurant", middleware.RequireAuth)
	{
		restaurantsRoutes.GET("/", restaurantsHandler.MyRestaurant)
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
