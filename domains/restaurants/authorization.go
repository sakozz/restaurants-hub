package restaurants

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
	ManagerId int64
}

func NewAuthorizer(currentUser *structs.BaseUser, managerId ...int64) Authorizor {
	if managerId == nil {
		managerId[0] = -1
	}
	return &authUser{currentUser, managerId[0]}
}

func (auth *authUser) AuthorizeAccess() bool {
	return auth.IsAdmin() || (auth.IsManager() && auth.UserOwnsResource())
}

func (auth *authUser) AuthorizeUpdate() bool {
	return auth.IsAdmin()
}

func (auth *authUser) AuthorizeDelete() bool {
	return auth.IsAdmin()
}

func (auth *authUser) UserOwnsResource() bool {
	return auth.Id == auth.ManagerId
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
		hasPermission = auth.Can("accessCollection", consts.Restaurants)
	case "create":
		hasPermission = auth.Can("create", consts.Restaurants)
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
