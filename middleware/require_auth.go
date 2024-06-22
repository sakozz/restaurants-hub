package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"resturants-hub.com/m/v2/domains/sessions"
	"resturants-hub.com/m/v2/domains/users"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

// MARK: RequireAuth
func RequireAuth(c *gin.Context) {
	/* Get cookie from request */
	tokenString, err := c.Cookie(os.Getenv("AUTH_COOKIE_NAME"))
	if tokenString == "" || err != nil {
		unauthorisedError(c)
		return
	}

	/* Validate  */
	sessionService := sessions.NewSessionService()
	session, restErr := sessionService.ValidateSessionToken(tokenString)
	if restErr != nil {
		c.AbortWithStatusJSON(restErr.Status(), restErr)
		return
	}

	// Get user by id
	usersDao := users.NewUsersDao()
	user, err := usersDao.GetSessionUser(&session.UserId)
	if err != nil {
		unauthorisedError(c)
		return
	}

	// set current user to be accessible for controllers
	c.Set("currentUser", user)
	c.Set("currentSession", session)
	c.Next()
}

// MARK: unauthorisedError
func unauthorisedError(c *gin.Context) {
	restErr := rest_errors.NewUnauthorizedError("Unauthorised Error")
	c.AbortWithStatusJSON(restErr.Status(), restErr)
}
