package services

import (
	"net/url"

	"resturants-hub.com/m/v2/domains/users"
	rest_errors "resturants-hub.com/m/v2/utils"
)

// var (
// 	UsersService usersServiceInterface = &usersService{}
// )

type UsersService interface {
	GetUser(int64) (*users.User, rest_errors.RestErr)
	SearchUser(url.Values) (users.Users, rest_errors.RestErr)
	UpdateUser(*users.User, interface{}) (*users.User, rest_errors.RestErr)
	// DeleteUser(int64) rest_errors.RestErr

	// LoginUser(users.LoginRequest) (*users.User, rest_errors.RestErr)
}

type usersService struct {
	dao users.UsersDao
}

func NewUsersService() UsersService {
	return &usersService{
		dao: users.NewUserDao(),
	}
}

func (service *usersService) GetUser(userId int64) (*users.User, rest_errors.RestErr) {

	result, err := service.dao.Get(&userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (service *usersService) UpdateUser(user *users.User, payload interface{}) (*users.User, rest_errors.RestErr) {

	updatedUser, err := service.dao.Update(user, payload)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// func (s *usersService) DeleteUser(userId int64) rest_errors.RestErr {
// 	dao := &users.User{Id: userId}
// 	return dao.Delete()
// }

func (s *usersService) SearchUser(params url.Values) (users.Users, rest_errors.RestErr) {
	return s.dao.Search(params)
}
