package users

import (
	consts "resturants-hub.com/m/v2/packages/const"
	"resturants-hub.com/m/v2/packages/structs"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type Authorizor interface {
	Authorize(string) (interface{}, rest_errors.RestErr)
	AuthorizeAccess() bool
	AuthorizeUpdate() bool
	AuthorizeDelete() bool
	userIsSelf() bool
}

type permissions struct {
	CanAccess bool `json:"canAccess"`
	CanUpdate bool `json:"canUpdate"`
	CanDelete bool `json:"canDelete"`
}

type authUser struct {
	*structs.BaseUser
	UserId int64
}

func NewAuthorizer(currentUser *structs.BaseUser, userId ...int64) Authorizor {
	if userId == nil {
		userId = append(userId, -1)
	}
	return &authUser{currentUser, userId[0]}
}

func (auth *authUser) AuthorizeAccess() bool {
	return auth.IsAdmin() || (auth.IsManager() && auth.userIsSelf())
}

func (auth *authUser) AuthorizeUpdate() bool {
	return auth.IsAdmin()
}

func (auth *authUser) AuthorizeDelete() bool {
	return auth.IsAdmin()
}

func (auth *authUser) userIsSelf() bool {
	return auth.Id == auth.UserId
}

func (auth *authUser) Authorize(action string) (interface{}, rest_errors.RestErr) {

	permissions := &permissions{
		CanAccess: auth.AuthorizeAccess(),
		CanUpdate: auth.AuthorizeUpdate(),
		CanDelete: auth.AuthorizeDelete(),
	}

	var hasPermission bool
	switch action {
	case "accessCollection":
		hasPermission = auth.Can("accessCollection", consts.Users)
	case "create":
		hasPermission = auth.Can("create", consts.Users)
	case "access":
		hasPermission = permissions.CanAccess
	case "update":
		hasPermission = permissions.CanUpdate
	case "delete":
		hasPermission = permissions.CanDelete
	default:
		hasPermission = false
	}

	if hasPermission {
		return permissions, nil
	}

	return nil, rest_errors.NewForbiddenError("You are not allowed to perform this action")
}
