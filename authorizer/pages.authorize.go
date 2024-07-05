package authorizer

import (
	"resturants-hub.com/m/v2/dto"
	consts "resturants-hub.com/m/v2/packages/const"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type PageAuthorizor interface {
	Authorize(string) (interface{}, rest_errors.RestErr)
	AuthorizeAccess() bool
	AuthorizeUpdate() bool
	AuthorizeDelete() bool
	UserOwnsResource() bool
}

type pagesAuthUser struct {
	*dto.BaseUser
	AuthorId int64
}

func NewPageAuthorizer(currentUser *dto.BaseUser, authorId ...int64) PageAuthorizor {
	if authorId == nil {
		authorId = []int64{0}
	}
	return &pagesAuthUser{currentUser, authorId[0]}
}

func (auth *pagesAuthUser) AuthorizeAccess() bool {
	return auth.IsAdmin() || (auth.IsManager() && auth.UserOwnsResource())
}

func (auth *pagesAuthUser) AuthorizeUpdate() bool {
	return auth.IsAdmin() || (auth.IsManager() && auth.UserOwnsResource())
}

func (auth *pagesAuthUser) AuthorizeDelete() bool {
	return auth.IsAdmin() || (auth.IsManager() && auth.UserOwnsResource())
}

func (auth *pagesAuthUser) UserOwnsResource() bool {
	return auth.Id == auth.AuthorId
}

/*
Use this pagePermissions and authorization for member resource only
The idea is to authorize the user based on the action they want to perform on the resource.
Authorization for collection/list of resources is handled in the handler via the AuthorizeCollection method
*/
type pagePermissions struct {
	CanAccess bool `json:"canAccess"`
	CanUpdate bool `json:"canUpdate"`
	CanDelete bool `json:"canDelete"`
}

func (auth *pagesAuthUser) Authorize(action string) (interface{}, rest_errors.RestErr) {
	permissions := &pagePermissions{
		CanAccess: auth.AuthorizeAccess(),
		CanUpdate: auth.AuthorizeUpdate(),
		CanDelete: auth.AuthorizeDelete(),
	}

	var hasPermission bool
	switch action {
	case "accessCollection":
		hasPermission = auth.Can("accessCollection", consts.Pages)
	case "create":
		hasPermission = auth.Can("create", consts.Pages)
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
