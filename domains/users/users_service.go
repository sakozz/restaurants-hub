package users

import (
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

// var (
// 	UsersService usersServiceInterface = &usersService{}
// )

type UsersService interface {
	GetUser(int64) (*User, rest_errors.RestErr)
	UpdateUser(*User, interface{}) (*User, rest_errors.RestErr)
	// DeleteUser(int64) rest_errors.RestErr

	// LoginUser(LoginRequest) (*User, rest_errors.RestErr)
}

type usersService struct {
	dao UsersDao
}

func NewUsersService() UsersService {
	return &usersService{
		dao: NewUserDao(),
	}
}

func (service *usersService) GetUser(userId int64) (*User, rest_errors.RestErr) {

	result, err := service.dao.Get(&userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (service *usersService) UpdateUser(user *User, payload interface{}) (*User, rest_errors.RestErr) {

	updatedUser, err := service.dao.Update(user, payload)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
