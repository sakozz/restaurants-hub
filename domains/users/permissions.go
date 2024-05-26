package users

import (
	"slices"

	consts "resturants-hub.com/m/v2/packages/const"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type PermissionsMap map[consts.ResourceType]interface{}

var (
	PermissionMappings PermissionsMap = PermissionsMap{
		consts.Restaurants: map[string][]consts.Role{
			"accessCollection": {consts.Admin, consts.Manager},
			"access":           {consts.Admin, consts.Manager},
			"create":           {consts.Admin},
			"update":           {consts.Admin, consts.Manager},
			"delete":           {consts.Admin},
		},
	}
)

func (user *User) Can(action string, resource consts.ResourceType) (bool, rest_errors.RestErr) {
	mappings := PermissionMappings[resource].(map[string][]consts.Role)
	roles := mappings[action]
	isPermitted := slices.Contains(roles, user.Role)

	if isPermitted {
		return isPermitted, nil
	}

	return isPermitted, rest_errors.NewForbiddenError("You are not allowed to perform this action")

}
