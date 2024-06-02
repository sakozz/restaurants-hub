package users

import (
	"slices"

	consts "resturants-hub.com/m/v2/packages/const"
)

type PermissionsMap map[consts.Role]interface{}

var (
	PermissionMappings PermissionsMap = PermissionsMap{
		consts.Admin: map[consts.ResourceType][]string{
			consts.Restaurants: {"accessCollection", "accessMember", "create"},
			consts.Users:       {"accessCollection", "accessMember", "create"},
			consts.Invitations: {"accessCollection", "accessMember", "create"},
		},
		consts.Manager: map[consts.ResourceType][]string{
			consts.Restaurants: {"accessMember"},
			consts.Users:       {"accessMember"},
			consts.Invitations: {},
		},
		consts.Public: map[consts.ResourceType][]string{
			consts.Restaurants: {},
			consts.Users:       {},
			consts.Invitations: {},
		},
	}
)

func (user *User) Can(action string, resource consts.ResourceType) bool {
	mappings := PermissionMappings[user.Role].(map[consts.ResourceType][]string)
	permissions := mappings[resource]
	return slices.Contains(permissions, action)
}

func (user *User) Permissions() interface{} {
	return PermissionMappings[user.Role]
}
