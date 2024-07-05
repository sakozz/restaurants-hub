package dto

import (
	"encoding/json"
	"strings"
	"time"

	consts "resturants-hub.com/m/v2/packages/const"
	"resturants-hub.com/m/v2/packages/types"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
	"resturants-hub.com/m/v2/serializers"
)

type User struct {
	BaseUser
	CreatedAt time.Time      `json:"createdAt" db:"created_at" goqu:"skipinsert"`
	UpdatedAt time.Time      `json:"updatedAt" db:"updated_at" goqu:"skipinsert"`
	DeletedAt types.NullTime `json:"deletedAt" db:"deleted_at" goqu:"skipinsert"`
}

type CreateUserPayload BaseUser

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PublicUserListItem struct {
	Id        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	AvatarURL string `json:"avatarUrl"`
}
type AdminUserListItem struct {
	PublicUserListItem
	Role  consts.Role `json:"role"`
	Email string      `json:"email"`
}

type AdminUserDetailItem struct {
	AdminUserListItem
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt types.NullTime `json:"deletedAt"`
}
type OwnerUserDetailItem struct {
	PublicUserListItem
	Email string      `json:"email"`
	Role  consts.Role `json:"role"`
}

type Users []User

func (user *User) UpdableAttributes() []string {
	return []string{"email", "username"}
}

func (user *User) Validate() rest_errors.RestErr {

	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return rest_errors.NewBadRequestError("invalid email address")
	}
	return nil
}

func (user *User) MemberFor(role consts.Role) interface{} {
	payload, _ := json.Marshal(user)
	switch role {

	case consts.Admin:
		var details AdminUserDetailItem
		json.Unmarshal(payload, &details)
		return serializers.MemberPayload[AdminUserDetailItem]{Id: user.Id, Type: "users", Attributes: details}
	default:
		var details OwnerUserDetailItem
		json.Unmarshal(payload, &details)
		return serializers.MemberPayload[OwnerUserDetailItem]{Id: user.Id, Type: "users", Attributes: details}
	}
}

func (users Users) CollectionFor(role consts.Role) []interface{} {
	result := make([]interface{}, len(users))
	for index, record := range users {
		payload, _ := json.Marshal(record)
		switch role {
		case consts.Admin:
			var adminListItem AdminUserListItem
			json.Unmarshal(payload, &adminListItem)
			result[index] = serializers.MemberPayload[AdminUserListItem]{Id: record.Id, Type: "users", Attributes: adminListItem}
		default:
			var publicListItem PublicUserListItem
			json.Unmarshal(payload, &publicListItem)
			result[index] = serializers.MemberPayload[PublicUserListItem]{Id: record.Id, Type: "users", Attributes: publicListItem}
		}
	}
	return result
}
