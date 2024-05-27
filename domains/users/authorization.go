package users

import (
	"github.com/gin-gonic/gin"
	consts "resturants-hub.com/m/v2/packages/const"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type authorizer struct {
	currentUser *User
	resource    *User
}

func (auth *authorizer) AuthorizeAccess() bool {
	return auth.currentUser.IsAdmin() || (auth.currentUser.IsManager() && auth.userIsSelf())
}

func (auth *authorizer) AuthorizeUpdate() bool {
	return auth.currentUser.IsAdmin()
}

func (auth *authorizer) AuthorizeDelete() bool {
	return auth.currentUser.IsAdmin()
}

func (auth *authorizer) userIsSelf() bool {
	if auth.resource == nil || auth.currentUser == nil {
		return false
	}
	return auth.currentUser.Id == auth.resource.Id
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

func (ctr *usersHandler) Authorize(action string, resource *User, c *gin.Context) (interface{}, rest_errors.RestErr) {

	// Get current user from context
	userData, ok := c.Get("currentUser")
	if !ok {
		return nil, rest_errors.NewUnauthorizedError("unauthorized")
	}

	ctr.currentUser = userData.(*User)

	authorizer := &authorizer{
		currentUser: ctr.currentUser,
		resource:    resource,
	}

	permissions := &permissions{
		CanAccess: authorizer.AuthorizeAccess(),
		CanUpdate: authorizer.AuthorizeUpdate(),
		CanDelete: authorizer.AuthorizeDelete(),
	}

	var hasPermission bool
	switch action {
	case "accessCollection":
		hasPermission = ctr.currentUser.Can("accessCollection", consts.Users)
	case "create":
		hasPermission = ctr.currentUser.Can("create", consts.Users)
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
