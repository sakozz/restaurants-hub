package users

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/jmoiron/sqlx"
	"resturants-hub.com/m/v2/database"
	rest_errors "resturants-hub.com/m/v2/utils"
)

type UsersDao interface {
	Create(*User) (*User, rest_errors.RestErr)
	FindOrCreate(*User) (*User, rest_errors.RestErr)
	Update(*User, interface{}) (*User, rest_errors.RestErr)
	Get(id *int64) (*User, rest_errors.RestErr)
	Search(url.Values) (Users, rest_errors.RestErr)
	Where(params map[string]interface{}) (*User, rest_errors.RestErr)
}

type connection struct {
	db         *sqlx.DB
	sqlBuilder database.SqlBuilder
}

func NewUserDao() UsersDao {
	return &connection{
		db:         database.DB,
		sqlBuilder: database.NewSqlBuilder(),
	}
}

const (
	Admin     AuthType = 0
	OwnerUser          = 1
	Public             = 2
)

func (user *User) UpdableAttributes() []string {
	return []string{"email", "username"}
}

func (user *User) Validate() rest_errors.RestErr {

	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return rest_errors.NewBadRequestError("invalid email address")
	}
	return nil
}

func (connection *connection) Create(payload *User) (*User, rest_errors.RestErr) {
	user := &User{}
	sqlQuery := connection.sqlBuilder.Insert("users", payload)

	row := connection.db.QueryRowx(sqlQuery)
	if row.Err() != nil {
		if uniquenessViolation, constraintName := database.HasUniquenessViolation(row.Err()); uniquenessViolation {
			return nil, rest_errors.InvalidError(ErrorMessage(constraintName))
		}
		return nil, rest_errors.NewInternalServerError("Server Error", row.Err())
	}

	row.StructScan(user)
	return user, nil
}

func (connection *connection) FindOrCreate(userData *User) (*User, rest_errors.RestErr) {

	user, _ := connection.Where(map[string]interface{}{
		"email": userData.Email,
	})

	if user != nil {
		return user, nil
	}

	sqlQuery := connection.sqlBuilder.Insert("users", userData)

	row := connection.db.QueryRowx(sqlQuery)
	if row.Err() != nil {
		if uniquenessViolation, constraintName := database.HasUniquenessViolation(row.Err()); uniquenessViolation {
			return nil, rest_errors.InvalidError(ErrorMessage(constraintName))
		}
		return nil, rest_errors.NewInternalServerError("Server Error", row.Err())
	}

	newUser := &User{}
	row.StructScan(newUser)

	return newUser, nil
}

func (connection *connection) Update(user *User, payload interface{}) (*User, rest_errors.RestErr) {
	sqlQuery := connection.sqlBuilder.Update("users", &user.ID, payload)
	row := connection.db.QueryRowx(sqlQuery)
	if row.Err() != nil {
		if uniquenessViolation, constraintName := database.HasUniquenessViolation(row.Err()); uniquenessViolation {
			return nil, rest_errors.InvalidError(ErrorMessage(constraintName))
		}
		return nil, rest_errors.NewInternalServerError("Server Error", row.Err())
	}
	row.StructScan(user)
	return user, nil
}

func (connection *connection) Get(id *int64) (*User, rest_errors.RestErr) {
	user := &User{}
	query := connection.sqlBuilder.Find("users", map[string]interface{}{"id": id})
	err := connection.db.Get(user, query)

	if err != nil {
		message := fmt.Sprintf("Sorry, user with id %v doesn't exist", *id)
		return nil, rest_errors.NewNotFoundError(message)
	}

	return user, nil
}

func (connection *connection) Where(params map[string]interface{}) (*User, rest_errors.RestErr) {
	user := &User{}

	query := connection.sqlBuilder.SearchBy("users", params)

	err := connection.db.Get(user, query)
	if err != nil {
		message := fmt.Sprintf("Sorry, user doesn't exist")
		return nil, rest_errors.NewNotFoundError(message)
	}

	return user, nil
}

func (connection *connection) Search(params url.Values) (Users, rest_errors.RestErr) {
	var users Users
	sqlQuery := connection.sqlBuilder.Filter("users", params)
	err := connection.db.Select(&users, sqlQuery)
	if err != nil {
		return nil, rest_errors.NewNotFoundError(err.Error())
	}

	return users, nil
}
