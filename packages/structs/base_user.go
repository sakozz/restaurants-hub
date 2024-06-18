package structs

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
			consts.Restaurants: {"accessMember", "create"},
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

type BaseUser struct {
	Id        int64       `json:"id" db:"id" goqu:"skipinsert"`
	Email     string      `json:"email" db:"email"`
	Role      consts.Role `json:"role" goqu:"skipinsert"`
	FirstName string      `json:"firstName" db:"first_name"`
	LastName  string      `json:"lastName" db:"last_name"`
	AvatarURL string      `json:"avatarUrl" db:"avatar_url"`
}

func (user *BaseUser) Can(action string, resource consts.ResourceType) bool {
	mappings := PermissionMappings[user.Role].(map[consts.ResourceType][]string)
	permissions := mappings[resource]
	return slices.Contains(permissions, action)
}

func (user *BaseUser) Permissions() interface{} {
	return PermissionMappings[user.Role]
}

func (user *BaseUser) IsAdmin() bool {
	return user.Role == consts.Admin
}

func (user *BaseUser) IsManager() bool {
	return user.Role == consts.Manager
}
