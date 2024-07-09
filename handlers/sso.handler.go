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
	"resturants-hub.com/m/v2/configs"
	"resturants-hub.com/m/v2/dao"
	"resturants-hub.com/m/v2/dto"
	consts "resturants-hub.com/m/v2/packages/const"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
	"resturants-hub.com/m/v2/services"
)

type SsoHandler interface {
	SsoLogin(c *gin.Context)
	RenewSession(c *gin.Context)
	Callback(c *gin.Context)
	Logout(c *gin.Context)
}

type ssoHandler struct {
	service        services.SessionService
	usersDao       dao.UsersDao
	invitationsDao dao.InvitationsDao
}

func NewSsoHandler() SsoHandler {
	return &ssoHandler{
		service:        services.NewSessionService(),
		usersDao:       dao.NewUsersDao(),
		invitationsDao: dao.NewInvitationDao(),
	}
}

/**
* SsoLogin handles the initial login process for single sign-on (SSO).
* This function generates a PKCE verifier and state, sets them in memory cache,
* and redirects the user to the SSO provider's authorization URL. The SSO provider is determined
* by the 'provider' query parameter from the request URL. If no provider is specified, it returns
* a bad request error with an appropriate message.
*
* @param c *gin.Context: The Gin context for handling HTTP requests and responses.
 */
func (handler *ssoHandler) SsoLogin(c *gin.Context) {
	// use PKCE to protect against CSRF attacks
	// https://www.ietf.org/archive/id/draft-ietf-oauth-security-topics-22.html#name-countermeasures-6
	verifier := oauth2.GenerateVerifier()
	state := "state"

	configs.MemoryCache.Set(state, verifier)

	// extract provider from URL and get new SSO config for provider
	provider := GetIdentifierFromUrl(c, "provider", false)
	if provider == "" {
		slugErr := rest_errors.NewBadRequestError("slug is required")
		c.JSON(slugErr.Status(), slugErr)
		return
	}
	ssoConfig := configs.NewSsoConfig(consts.SsoProvider(provider))

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := ssoConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	c.Redirect(http.StatusFound, url)
}

func (handler *ssoHandler) Callback(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	ctx := context.Background()
	code := queryParams.Get("code")
	state := queryParams.Get("state")

	verifier, err := configs.MemoryCache.Get(state)
	if err != nil {
		fmt.Println(err)
	}

	// extract provider from URL and get new SSO config for provider
	provider := GetIdentifierFromUrl(c, "provider", false)
	if provider == "" {
		slugErr := rest_errors.NewBadRequestError("slug is required")
		c.JSON(slugErr.Status(), slugErr)
		return
	}
	ssoConfig := configs.NewSsoConfig(consts.SsoProvider(provider))

	token, err := ssoConfig.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		fmt.Println("Error on SSO token exchange:", err)
		restErr := rest_errors.NewInternalServerError(err)
		c.JSON(restErr.Status(), restErr)
		return
	}

	// Retrieve user data by token
	client := ssoConfig.Client(c, token)
	userInfourl := ssoConfig.UserInfoUrl + "?access_token=" + token.AccessToken
	userData, err := handler.RetrieveUserInfo(client, userInfourl)
	if err != nil {
		fmt.Println("Retrieve user error:", err)
		restErr := rest_errors.NewInternalServerError(err)
		c.JSON(restErr.Status(), restErr)
		return
	}

	// Check if user is registered
	user := handler.usersDao.Where(map[string]interface{}{"email": userData.Email})
	if user == nil {
		// Check if user has a valid invitation

		invitation := handler.validateInvitation(userData.Email, c)

		// If user is not registered, check if user has a valid invitation
		if invitation == nil {
			return
		}

		// Create new user with role from invitation
		userData.Role = invitation.Role
		newUser, restErr := handler.usersDao.CreateUser(userData)
		if restErr != nil {
			fmt.Println("New user created:", restErr)
			c.JSON(restErr.Status(), restErr)
			return
		}
		user = newUser
	}

	// save session
	session := (dto.Session{
		Provider:     "google",
		AccessToken:  token.AccessToken,
		ExpiresAt:    token.Expiry,
		RefreshToken: token.RefreshToken,
		Email:        user.Email,
		UserId:       user.Id,
	})
	_, error := handler.service.CreateSession(&session)

	// save session and return user
	if error != nil {
		c.JSON(error.Status(), err)
		return
	}
	// Finally, we set the client cookie for "token"
	setCookie(c, &session)

	c.JSON(http.StatusOK, session)
}

func (handler *ssoHandler) RenewSession(c *gin.Context) {
	currentSession, sessionErr := c.Get("currentSession")

	if !sessionErr {
		restErr := rest_errors.NewUnauthorizedError("Unauthorised user. No active session")
		c.JSON(restErr.Status(), restErr)
		return
	}
	session := currentSession.(*dto.Session)

	// extract provider from URL and get new SSO config for provider
	provider := GetIdentifierFromUrl(c, "provider", false)
	if provider == "" {
		slugErr := rest_errors.NewBadRequestError("slug is required")
		c.JSON(slugErr.Status(), slugErr)
		return
	}
	ssoConfig := configs.NewSsoConfig(consts.SsoProvider(provider))

	tokenSource := ssoConfig.TokenSource(context.Background(), &oauth2.Token{
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
	newSession, restErr := handler.service.RenewSession(&dto.Session{
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

	// Finally, we set the client cookie for "token"
	setCookie(c, newSession)

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
	session := currentSession.(*dto.Session)
	_, err := handler.service.InvalidateToken(session)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	// Clear cookie and return success response
	c.SetCookie(os.Getenv("AUTH_COOKIE_NAME"), "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, "Success")
}

func (handler *ssoHandler) RetrieveUserInfo(client *http.Client, userInfourl string) (*dto.CreateUserPayload, error) {
	response, err := client.Get(userInfourl)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	data := &dto.SsoUserInfo{}
	err = json.Unmarshal(contents, data)
	if err != nil {
		panic(err)
	}

	return &dto.CreateUserPayload{
		Email:     data.Email,
		FirstName: data.GivenName,
		LastName:  data.FamilyName,
		AvatarURL: data.Picture,
	}, nil
}

func setCookie(c *gin.Context, session *dto.Session) {
	c.SetCookie(os.Getenv("AUTH_COOKIE_NAME"), session.AccessToken, 2000, "/", "localhost", false, false)
}

func (handler *ssoHandler) validateInvitation(email string, c *gin.Context) *dto.Invitation {
	// Check if user has a valid invitation
	invitation := handler.invitationsDao.SearchInvitations(map[string]interface{}{"email": email})

	invitationErr := rest_errors.NewForbiddenError("User is not registered or no valid invitation")
	// If user is not registered, check if user has a valid invitation
	if invitation == nil || !invitation.IsValid() {
		c.JSON(invitationErr.Status(), invitationErr)
		return nil
	}

	if !invitation.IsValid() {
		c.JSON(invitationErr.Status(), invitationErr)
		return nil
	}

	return invitation
}
