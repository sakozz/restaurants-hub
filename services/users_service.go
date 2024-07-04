package services

import (
	"resturants-hub.com/m/v2/dao"
	"resturants-hub.com/m/v2/dto"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

// var (
// 	UsersService usersServiceInterface = &usersService{}
// )

type UsersService interface {
	GetUser(int64) (*dto.User, rest_errors.RestErr)
	UpdateUser(*dto.User, interface{}) (*dto.User, rest_errors.RestErr)
}

type usersService struct {
	dao dao.UsersDao
}

func NewUsersService() UsersService {
	return &usersService{
		dao: dao.NewUsersDao(),
	}
}

func (service *usersService) GetUser(userId int64) (*dto.User, rest_errors.RestErr) {

	result, err := service.dao.GetUser(&userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (service *usersService) UpdateUser(user *dto.User, payload interface{}) (*dto.User, rest_errors.RestErr) {

	updatedUser, err := service.dao.UpdateUser(&user.Id, payload)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
