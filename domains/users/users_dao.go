package users

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/jmoiron/sqlx"
	"resturants-hub.com/m/v2/database"
	consts "resturants-hub.com/m/v2/packages/const"
	"resturants-hub.com/m/v2/packages/structs"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

// MARK: UsersDao
type UsersDao interface {
	Create(*CreateUserPayload) (*User, rest_errors.RestErr)
	FindOrCreate(*CreateUserPayload) (*User, rest_errors.RestErr)
	Update(*User, interface{}) (*User, rest_errors.RestErr)
	Get(id *int64) (*User, rest_errors.RestErr)
	GetSessionUser(id *int64) (*structs.BaseUser, rest_errors.RestErr)
	AuthorizedCollection(url.Values, *structs.BaseUser) (Users, rest_errors.RestErr)
	Where(params map[string]interface{}) *User
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

func (connection *connection) Create(payload *CreateUserPayload) (*User, rest_errors.RestErr) {
	user := &User{}
	sqlQuery := connection.sqlBuilder.Insert("users", payload)

	row := connection.db.QueryRowx(sqlQuery)
	if row.Err() != nil {
		if uniquenessViolation, constraintName := database.HasUniquenessViolation(row.Err()); uniquenessViolation {
			return nil, rest_errors.InvalidError(ErrorMessage(constraintName))
		}
		return nil, rest_errors.NewInternalServerError(row.Err())
	}

	row.StructScan(user)
	return user, nil
}

func (connection *connection) FindOrCreate(userData *CreateUserPayload) (*User, rest_errors.RestErr) {

	user := connection.Where(map[string]interface{}{
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
		return nil, rest_errors.NewInternalServerError(row.Err())
	}

	newUser := &User{}
	row.StructScan(newUser)

	return newUser, nil
}

func (connection *connection) Update(user *User, payload interface{}) (*User, rest_errors.RestErr) {
	sqlQuery := connection.sqlBuilder.Update("users", &user.Id, payload)
	row := connection.db.QueryRowx(sqlQuery)
	if row.Err() != nil {
		if uniquenessViolation, constraintName := database.HasUniquenessViolation(row.Err()); uniquenessViolation {
			return nil, rest_errors.InvalidError(ErrorMessage(constraintName))
		}
		return nil, rest_errors.NewInternalServerError(row.Err())
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

func (connection *connection) GetSessionUser(id *int64) (*structs.BaseUser, rest_errors.RestErr) {
	user := &User{}
	query := connection.sqlBuilder.Find("users", map[string]interface{}{"id": id})
	err := connection.db.Get(user, query)

	if err != nil {
		message := fmt.Sprintf("Sorry, user with id %v doesn't exist", *id)
		return nil, rest_errors.NewNotFoundError(message)
	}

	return &user.BaseUser, nil
}

func (connection *connection) Where(params map[string]interface{}) *User {
	user := &User{}

	query := connection.sqlBuilder.SearchBy("users", params)
	err := connection.db.Get(user, query)
	if err != nil {
		fmt.Println("Error Occured:", err)
		return nil
	}

	return user
}

func (connection *connection) AuthorizedCollection(params url.Values, user *structs.BaseUser) (Users, rest_errors.RestErr) {
	switch user.Role {
	case consts.Admin:
		return connection.search(params)
	default:
		return Users{}, nil
	}
}

func (connection *connection) search(params url.Values) (Users, rest_errors.RestErr) {
	var users Users
	sqlQuery := connection.sqlBuilder.Filter("users", params)
	err := connection.db.Select(&users, sqlQuery)
	if err != nil {
		return nil, rest_errors.NewNotFoundError(err.Error())
	}

	return users, nil
}
