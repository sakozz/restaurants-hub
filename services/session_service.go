package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"resturants-hub.com/m/v2/domains/users"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

var JwtLifeSpan int = 5 * 60 // 5 minutes

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Id int64 `json:"id"`
	jwt.RegisteredClaims
}

type Jwt struct {
	MaxAge         int
	ExpirationTime time.Time
	Token          string
}

type SessionService interface {
	CreateSession(*users.Session) (*users.Session, rest_errors.RestErr)
	InvalidateToken(*users.Session) (bool, rest_errors.RestErr)
	ValidateSessionToken(string) (*users.Session, rest_errors.RestErr)
	RenewSession(*users.Session) (*users.Session, rest_errors.RestErr)
	/*
	   FindUserByLoginPayload(users.Session) (*users.Session, rest_errors.RestErr)
	   GenerateJwtToken(*users.Session) (*Jwt, error)
	   RenewSession(*users.User, *Jwt) (bool, rest_errors.RestErr)
	   ValidateSessionToken(string) (*tokens.Token, rest_errors.RestErr)

	*/
}

type sessionService struct {
	sessionDao users.SessionDao
	usersDao   users.UsersDao
}

func NewSessionService() SessionService {
	return &sessionService{
		sessionDao: users.NewSessionDao(),
		usersDao:   users.NewUserDao(),
	}
}

func (service *sessionService) CreateSession(userSession *users.Session) (*users.Session, rest_errors.RestErr) {
	// copier.Copy(userSession, gothUser)
	session, sessionError := service.sessionDao.Create(userSession)
	if sessionError != nil {
		fmt.Println(sessionError)
		return nil, sessionError
	}
	return session, nil
}

func (service *sessionService) RenewSession(session *users.Session) (*users.Session, rest_errors.RestErr) {
	fmt.Println(session.Id, session.ExpiresAt, session.AccessToken)
	session, tokenError := service.sessionDao.Update(session, &session.Id)
	if tokenError != nil {
		return nil, tokenError
	}
	return session, nil
}

/* func (service *sessionService) FindUserByLoginPayload(payload users.LoginPayload) (*users.User, rest_errors.RestErr) {
	return service.usersDao.FindByLoginPayload(payload)
}  */

func (service *sessionService) GenerateJwtToken(user *users.Session) (*Jwt, error) {
	// Declare the expiration time of the token here
	expirationTime := time.Now().Add(time.Duration(JwtLifeSpan) * time.Second)

	// Create the JWT claims, which includes the username and expiry time
	claims := Claims{
		// Id: user.,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	// Create the JWT key used to create the signature
	var jwtKey = []byte(os.Getenv("JWT_SECRET"))

	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	return &Jwt{
		MaxAge:         JwtLifeSpan,
		ExpirationTime: expirationTime,
		Token:          tokenString,
	}, nil
}

func (service *sessionService) ValidateSessionToken(token string) (*users.Session, rest_errors.RestErr) {
	params := map[string]interface{}{
		"access_token": token,
	}
	sessionToken, err := service.sessionDao.Find(params)

	if err != nil {
		return nil, rest_errors.NewUnauthorizedError("Unauthorised Error")
	}

	if sessionToken.ExpiresAt.Before(time.Now()) {
		return nil, rest_errors.NewUnauthorizedError("Session expired")
	}

	return sessionToken, nil
}

func (service *sessionService) InvalidateToken(session *users.Session) (bool, rest_errors.RestErr) {
	_, sessionError := service.sessionDao.ExpireToken(session)
	if sessionError != nil {
		return false, sessionError
	}
	return true, nil
}
