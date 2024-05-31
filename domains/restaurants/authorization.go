package restaurants

import (
	"github.com/gin-gonic/gin"
	"resturants-hub.com/m/v2/domains/users"
	consts "resturants-hub.com/m/v2/packages/const"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type authorizer struct {
	currentUser *users.User
	restaurant  *Restaurant
}

func (auth *authorizer) AuthorizeAccess() bool {
	return auth.currentUser.IsAdmin() || (auth.currentUser.IsManager() && auth.userOwnsResource())
}

func (auth *authorizer) AuthorizeUpdate() bool {
	return auth.currentUser.IsAdmin()
}

func (auth *authorizer) AuthorizeDelete() bool {
	return auth.currentUser.IsAdmin()
}

func (auth *authorizer) userOwnsResource() bool {
	if auth.restaurant == nil || auth.currentUser == nil {
		return false
	}
	return auth.currentUser.Id == auth.restaurant.UserId
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

func (ctr *restaurantsHandler) Authorize(action string, resource *Restaurant, c *gin.Context) (interface{}, rest_errors.RestErr) {

	// Get current user from context
	userData, ok := c.Get("currentUser")
	if !ok {
		return nil, rest_errors.NewUnauthorizedError("unauthorized")
	}

	ctr.currentUser = userData.(*users.User)

	authorizer := &authorizer{
		currentUser: ctr.currentUser,
		restaurant:  resource,
	}

	permissions := &permissions{
		CanAccess: authorizer.AuthorizeAccess(),
		CanUpdate: authorizer.AuthorizeUpdate(),
		CanDelete: authorizer.AuthorizeDelete(),
	}

	var hasPermission bool
	switch action {
	case "accessCollection":
		hasPermission = ctr.currentUser.Can("accessCollection", consts.Restaurants)
	case "create":
		hasPermission = ctr.currentUser.Can("create", consts.Restaurants)
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
