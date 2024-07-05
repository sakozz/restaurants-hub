package authorizer

import (
	"resturants-hub.com/m/v2/dto"
	consts "resturants-hub.com/m/v2/packages/const"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type InvitationAuthorizor interface {
	Authorize(string) (interface{}, rest_errors.RestErr)
	AuthorizeAccess() bool
	AuthorizeUpdate() bool
	AuthorizeDelete() bool
	UserOwnsResource() bool
}

type invitationsAuthUser struct {
	*dto.BaseUser
	InvitedEmail string
}

func NewInvitationAuthorizer(currentUser *dto.BaseUser, email ...string) InvitationAuthorizor {
	if len(email) == 0 {
		email = append(email, "")
	}
	return &invitationsAuthUser{currentUser, email[0]}

}

func (auth *invitationsAuthUser) AuthorizeAccess() bool {
	return auth.IsAdmin()
}

func (auth *invitationsAuthUser) AuthorizeUpdate() bool {
	return auth.IsAdmin()
}

func (auth *invitationsAuthUser) AuthorizeDelete() bool {
	return auth.IsAdmin()
}

func (auth *invitationsAuthUser) UserOwnsResource() bool {
	return auth.InvitedEmail == auth.Email
}

/*
Use this permissions and authorization for member resource only
The idea is to authorize the user based on the action they want to perform on the resource.
Authorization for collection/list of resources is handled in the handler via the AuthorizeCollection method
*/
type invitationPermissions struct {
	CanAccess bool `json:"canAccess"`
	CanUpdate bool `json:"canUpdate"`
	CanDelete bool `json:"canDelete"`
}

func (auth *invitationsAuthUser) Authorize(action string) (interface{}, rest_errors.RestErr) {

	permissions := &invitationPermissions{
		CanAccess: auth.AuthorizeAccess(),
		CanUpdate: auth.AuthorizeUpdate(),
		CanDelete: auth.AuthorizeDelete(),
	}

	var hasPermission bool
	switch action {
	case "accessCollection":
		hasPermission = auth.Can("accessCollection", consts.Invitations)
	case "create":
		hasPermission = auth.Can("create", consts.Invitations)
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
