package users

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"resturants-hub.com/m/v2/database"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type Session struct {
	Id                int64     `json:"id" db:"id" goqu:"skipinsert"`
	UserId            int64     `json:"userId" db:"user_id" goqu:"omitempty"`
	Provider          string    `json:"provider" db:"provider" goqu:"omitempty"`
	Email             string    `json:"email" db:"email" goqu:"omitempty"`
	AccessToken       string    `json:"accessToken" db:"access_token" goqu:"omitempty"`
	AccessTokenSecret string    `json:"accessTokenSecret" db:"access_token_secret" goqu:"omitempty"`
	RefreshToken      string    `json:"refreshToken" db:"refresh_token" goqu:"omitempty"`
	ExpiresAt         time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at" goqu:"skipinsert omitempty"`
	UpdatedAt         time.Time `json:"updatedAt" db:"updated_at" goqu:"skipinsert"`
	IDToken           string    `json:"idToken" db:"id_token"`
}

type sessionConnection struct {
	db         *sqlx.DB
	sqlBuilder database.SqlBuilder
}

type SessionDao interface {
	Create(*Session) (*Session, rest_errors.RestErr)
	Find(map[string]interface{}) (*Session, rest_errors.RestErr)
	ExpireToken(*Session) (bool, rest_errors.RestErr)
	Update(*Session, *int64) (*Session, rest_errors.RestErr)
}

func NewSessionDao() SessionDao {
	return &sessionConnection{
		db:         database.DB,
		sqlBuilder: database.NewSqlBuilder(),
	}
}

func (connection *sessionConnection) Create(payload *Session) (*Session, rest_errors.RestErr) {
	session := &Session{}
	sqlQuery := connection.sqlBuilder.Insert("sessions", payload)

	row := connection.db.QueryRowx(sqlQuery)
	if row.Err() != nil {
		fmt.Println(row.Err())
		if uniquenessViolation, constraintName := database.HasUniquenessViolation(row.Err()); uniquenessViolation {
			return nil, rest_errors.InvalidError(ErrorMessage(constraintName))
		}
		return nil, rest_errors.NewInternalServerError(row.Err())
	}

	row.StructScan(session)
	return session, nil
}

func (connection *sessionConnection) Update(payload *Session, sessionId *int64) (*Session, rest_errors.RestErr) {
	session := &Session{}
	sqlQuery := connection.sqlBuilder.Update("sessions", sessionId, payload)

	row := connection.db.QueryRowx(sqlQuery)
	if row.Err() != nil {
		fmt.Println(row.Err())
		if uniquenessViolation, constraintName := database.HasUniquenessViolation(row.Err()); uniquenessViolation {
			return nil, rest_errors.InvalidError(ErrorMessage(constraintName))
		}
		return nil, rest_errors.NewInternalServerError(row.Err())
	}
	row.StructScan(session)
	return session, nil
}

func (connection *sessionConnection) Find(params map[string]interface{}) (*Session, rest_errors.RestErr) {
	session := &Session{}

	query := connection.sqlBuilder.SearchBy("sessions", params)

	err := connection.db.Get(session, query)
	if err != nil {
		fmt.Println("cookie: ", err)
		message := fmt.Sprintf("Failed to find token record for parameter %v", params)
		return nil, rest_errors.NewNotFoundError(message)
	}

	return session, nil
}

func (connection *sessionConnection) ExpireToken(session *Session) (bool, rest_errors.RestErr) {

	session.ExpiresAt = time.Now()
	query := connection.sqlBuilder.Update("sessions", &session.Id, session)

	row := connection.db.QueryRowx(query)
	if row.Err() != nil {
		if uniquenessViolation, constraintName := database.HasUniquenessViolation(row.Err()); uniquenessViolation {
			return true, rest_errors.InvalidError(ErrorMessage(constraintName))
		}
		return true, rest_errors.NewInternalServerError(row.Err())
	}

	return true, nil
}

func ErrorMessage(errorKey string) string {
	errors := map[string]string{
		"users_username_key": "Username must be unique",
		"users_email_key":    "Email must be unique",
		"user_not_found":     "User is not found",
	}

	return errors[errorKey]
}
