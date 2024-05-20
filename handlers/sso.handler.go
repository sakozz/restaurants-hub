package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"resturants-hub.com/m/v2/domains/sso"
	"resturants-hub.com/m/v2/domains/users"
	"resturants-hub.com/m/v2/services"
	rest_errors "resturants-hub.com/m/v2/utils"
)

type SsoHandler interface {
	SsoLogin(c *gin.Context)
	RenewSession(c *gin.Context)
	Callback(c *gin.Context)
	Logout(c *gin.Context)
}

type ssoHandler struct {
	service  services.SessionService
	usersDao users.UsersDao
	/*payload RequestPayload */
}

func NewSsoHandler() SsoHandler {
	return &ssoHandler{
		service:  services.NewSessionService(),
		usersDao: users.NewUserDao(),
	}
}

func (handler *ssoHandler) SsoLogin(c *gin.Context) {
	// use PKCE to protect against CSRF attacks
	// https://www.ietf.org/archive/id/draft-ietf-oauth-security-topics-22.html#name-countermeasures-6
	verifier := oauth2.GenerateVerifier()
	state := "state"

	sso.MemoryCache.Set(state, verifier)

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := sso.GoogleSsoConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	c.Redirect(http.StatusFound, url)
}

func (handler *ssoHandler) Callback(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	ctx := context.Background()
	code := queryParams.Get("code")
	state := queryParams.Get("state")

	verifier, err := sso.MemoryCache.Get(state)
	if err != nil {
		fmt.Println(err)
	}

	token, err := sso.GoogleSsoConfig.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		fmt.Println(err)
	}

	// Retrieve user data by token
	client := sso.GoogleSsoConfig.Client(c, token)
	userData, err := handler.RetrieveUserInfo(client, token.AccessToken)
	if err != nil {
		fmt.Println(err)
	}

	user, restErr := handler.usersDao.FindOrCreate(userData)
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
	}

	// save session
	session := (users.Session{
		Provider:     "google",
		AccessToken:  token.AccessToken,
		ExpiresAt:    token.Expiry,
		RefreshToken: token.RefreshToken,
		Email:        user.Email,
		ProfileId:    user.ID,
	})
	_, error := handler.service.CreateSession(&session)

	// save session and return user
	if error != nil {
		c.JSON(error.Status(), err)
		return
	}
	// Finally, we set the client cookie for "token"
	// we also set an expiry time which is the same as the token itself
	// c.SetCookie("restaurant-cookie", session.AccessToken, 2000, "/", "localhost", false, true)
	c.JSON(http.StatusOK, session)
}

func (handler *ssoHandler) RenewSession(c *gin.Context) {
	currentSession, sessionErr := c.Get("currentSession")

	if !sessionErr {
		restErr := rest_errors.NewUnauthorizedError("Unauthorised user. No active session")
		c.JSON(restErr.Status(), restErr)
		return
	}
	session := currentSession.(*users.Session)

	tokenSource := sso.GoogleSsoConfig.TokenSource(context.Background(), &oauth2.Token{
		RefreshToken: session.RefreshToken,
	})

	token, error := tokenSource.Token()
	if error != nil {
		restErr := rest_errors.NewUnauthorizedError("Failed to renew session")
		c.JSON(restErr.Status(), restErr)
		return
	}

	// fmt.Println(token.Expiry, token.AccessToken)

	// save session
	newSession, restErr := handler.service.RenewSession(&users.Session{
		Id:           session.Id,
		AccessToken:  token.AccessToken,
		ExpiresAt:    token.Expiry,
		RefreshToken: token.RefreshToken,
	})

	// save session and return user
	if restErr != nil {
		c.JSON(restErr.Status(), restErr)
		return
	}

	c.JSON(http.StatusOK, newSession)
}

func (handler *ssoHandler) Logout(c *gin.Context) {
	// Extract current session token from context
	currentSession, exists := c.Get("currentSession")

	if !exists {
		restErr := rest_errors.NewUnauthorizedError("Unauthorised user. No active session")
		c.JSON(restErr.Status(), restErr)
		return
	}

	// Invalidate token by setting expired to current DateTime.
	session := currentSession.(*users.Session)
	_, err := handler.service.InvalidateToken(session)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	// Clear cookie and return success response
	c.SetCookie(os.Getenv("AUTH_COOKIE_NAME"), "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, "Success")
}

func (handler *ssoHandler) RetrieveUserInfo(client *http.Client, token string) (*users.User, error) {
	userInfourl := os.Getenv("SSO_USER_INFO_URL") + "?access_token=" + token
	response, err := client.Get(userInfourl)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	data := &users.SsoUserInfo{}
	err = json.Unmarshal(contents, data)
	if err != nil {
		panic(err)
	}

	return &users.User{
		Email:     data.Email,
		FirstName: data.GivenName,
		LastName:  data.FamilyName,
		AvatarURL: data.Picture,
	}, nil
}
