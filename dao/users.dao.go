package dao

import (
	"fmt"
	"net/url"

	"resturants-hub.com/m/v2/database"
	"resturants-hub.com/m/v2/dto"
	consts "resturants-hub.com/m/v2/packages/const"
	"resturants-hub.com/m/v2/packages/structs"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type UsersDao interface {
	CreateUser(*dto.CreateUserPayload) (*dto.User, rest_errors.RestErr)
	FindOrCreateUser(*dto.CreateUserPayload) (*dto.User, rest_errors.RestErr)
	UpdateUser(id *int64, payload interface{}) (*dto.User, rest_errors.RestErr)
	GetUser(id *int64) (*dto.User, rest_errors.RestErr)
	GetSessionUser(id *int64) (*structs.BaseUser, rest_errors.RestErr)
	AuthorizedUsersCollection(url.Values, *structs.BaseUser) (dto.Users, rest_errors.RestErr)
	Where(params map[string]interface{}) *dto.User
}

func NewUsersDao() UsersDao {
	return &connection{
		db:         database.DB,
		sqlBuilder: database.NewSqlBuilder(),
	}
}

func (connection *connection) CreateUser(payload *dto.CreateUserPayload) (*dto.User, rest_errors.RestErr) {
	user := &dto.User{}
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

func (connection *connection) FindOrCreateUser(userData *dto.CreateUserPayload) (*dto.User, rest_errors.RestErr) {

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

	newUser := &dto.User{}
	row.StructScan(newUser)

	return newUser, nil
}

func (connection *connection) UpdateUser(id *int64, payload interface{}) (*dto.User, rest_errors.RestErr) {
	sqlQuery := connection.sqlBuilder.Update("users", id, payload)
	row := connection.db.QueryRowx(sqlQuery)
	if row.Err() != nil {
		if uniquenessViolation, constraintName := database.HasUniquenessViolation(row.Err()); uniquenessViolation {
			return nil, rest_errors.InvalidError(ErrorMessage(constraintName))
		}
		return nil, rest_errors.NewInternalServerError(row.Err())
	}
	user := &dto.User{}
	row.StructScan(user)
	return user, nil
}

func (connection *connection) GetUser(id *int64) (*dto.User, rest_errors.RestErr) {
	user := &dto.User{}
	query := connection.sqlBuilder.Find("users", map[string]interface{}{"id": id})
	err := connection.db.Get(user, query)

	if err != nil {
		message := fmt.Sprintf("Sorry, user with id %v doesn't exist", *id)
		return nil, rest_errors.NewNotFoundError(message)
	}

	return user, nil
}

func (connection *connection) GetSessionUser(id *int64) (*structs.BaseUser, rest_errors.RestErr) {
	user := &dto.User{}
	query := connection.sqlBuilder.Find("users", map[string]interface{}{"id": id})
	err := connection.db.Get(user, query)

	if err != nil {
		message := fmt.Sprintf("Sorry, user with id %v doesn't exist", *id)
		return nil, rest_errors.NewNotFoundError(message)
	}

	return &user.BaseUser, nil
}

func (connection *connection) Where(params map[string]interface{}) *dto.User {
	user := &dto.User{}

	query := connection.sqlBuilder.SearchBy("users", params)
	err := connection.db.Get(user, query)
	if err != nil {
		fmt.Println("Error Occured:", err)
		return nil
	}

	return user
}

func (connection *connection) AuthorizedUsersCollection(params url.Values, user *structs.BaseUser) (dto.Users, rest_errors.RestErr) {
	switch user.Role {
	case consts.Admin:
		return connection.searchUsers(params)
	default:
		return dto.Users{}, nil
	}
}

func (connection *connection) searchUsers(params url.Values) (dto.Users, rest_errors.RestErr) {
	var users dto.Users
	sqlQuery := connection.sqlBuilder.Filter("users", params)
	err := connection.db.Select(&users, sqlQuery)
	if err != nil {
		return nil, rest_errors.NewNotFoundError(err.Error())
	}

	return users, nil
}
