package dao

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"resturants-hub.com/m/v2/database"
	"resturants-hub.com/m/v2/dto"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type sessionConnection struct {
	db         *sqlx.DB
	sqlBuilder database.SqlBuilder
}

type SessionDao interface {
	CreateSession(*dto.Session) (*dto.Session, rest_errors.RestErr)
	FindSession(map[string]interface{}) (*dto.Session, rest_errors.RestErr)
	ExpireToken(*dto.Session) (bool, rest_errors.RestErr)
	UpdateSession(*dto.Session, *int64) (*dto.Session, rest_errors.RestErr)
}

func NewSessionDao() SessionDao {
	return &sessionConnection{
		db:         database.DB,
		sqlBuilder: database.NewSqlBuilder(),
	}
}

func (connection *sessionConnection) CreateSession(payload *dto.Session) (*dto.Session, rest_errors.RestErr) {
	session := &dto.Session{}
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

func (connection *sessionConnection) UpdateSession(payload *dto.Session, sessionId *int64) (*dto.Session, rest_errors.RestErr) {
	session := &dto.Session{}
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

func (connection *sessionConnection) FindSession(params map[string]interface{}) (*dto.Session, rest_errors.RestErr) {
	session := &dto.Session{}

	query := connection.sqlBuilder.SearchBy("sessions", params)

	err := connection.db.Get(session, query)
	if err != nil {
		fmt.Println("cookie: ", err)
		message := fmt.Sprintf("Failed to find token record for parameter %v", params)
		return nil, rest_errors.NewNotFoundError(message)
	}

	return session, nil
}

func (connection *sessionConnection) ExpireToken(session *dto.Session) (bool, rest_errors.RestErr) {

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
