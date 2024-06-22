package pages

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
	UserOwnsResource() bool
}
type authUser struct {
	*structs.BaseUser
	AuthorId int64
}

func NewAuthorizer(currentUser *structs.BaseUser, authorId ...int64) Authorizor {
	if authorId == nil {
		authorId = []int64{0}
	}
	return &authUser{currentUser, authorId[0]}
}

func (auth *authUser) AuthorizeAccess() bool {
	return auth.IsAdmin() || (auth.IsManager() && auth.UserOwnsResource())
}

func (auth *authUser) AuthorizeUpdate() bool {
	return auth.IsAdmin() || (auth.IsManager() && auth.UserOwnsResource())
}

func (auth *authUser) AuthorizeDelete() bool {
	return auth.IsAdmin() || (auth.IsManager() && auth.UserOwnsResource())
}

func (auth *authUser) UserOwnsResource() bool {
	return auth.Id == auth.AuthorId
}

/*
Use this permissions and authorization for member resource only
The idea is to authorize the user based on the action they want to perform on the resource.
Authorization for collection/list of resources is handled in the handler via the AuthorizeCollection method
*/
type permissions struct {
	CanAccess bool `json:"canAccess"`
	CanUpdate bool `json:"canUpdate"`
	CanDelete bool `json:"canDelete"`
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
