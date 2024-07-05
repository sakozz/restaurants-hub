package authorizer

import (
	"resturants-hub.com/m/v2/dto"
	consts "resturants-hub.com/m/v2/packages/const"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type RestaurantAuthorizor interface {
	Authorize(string) (interface{}, rest_errors.RestErr)
	AuthorizeAccess() bool
	AuthorizeUpdate() bool
	AuthorizeDelete() bool
	UserOwnsResource() bool
}
type restaurantAuthUser struct {
	*dto.BaseUser
	ManagerId int64
}

func NewRestaurantsAuthorizer(currentUser *dto.BaseUser, managerId ...int64) RestaurantAuthorizor {
	if managerId == nil {
		managerId = []int64{0}
	}
	return &restaurantAuthUser{currentUser, managerId[0]}
}

func (auth *restaurantAuthUser) AuthorizeAccess() bool {
	return auth.IsAdmin() || (auth.IsManager() && auth.UserOwnsResource())
}

func (auth *restaurantAuthUser) AuthorizeUpdate() bool {
	return auth.IsAdmin()
}

func (auth *restaurantAuthUser) AuthorizeDelete() bool {
	return auth.IsAdmin()
}

func (auth *restaurantAuthUser) UserOwnsResource() bool {
	return auth.Id == auth.ManagerId
}

/*
Use this permissions and authorization for member resource only
The idea is to authorize the user based on the action they want to perform on the resource.
Authorization for collection/list of resources is handled in the handler via the AuthorizeCollection method
*/
type restaurantPermissions struct {
	CanAccess bool `json:"canAccess"`
	CanUpdate bool `json:"canUpdate"`
	CanDelete bool `json:"canDelete"`
}

func (auth *restaurantAuthUser) Authorize(action string) (interface{}, rest_errors.RestErr) {
	permissions := &restaurantPermissions{
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
