package users

import (
	"time"

	consts "resturants-hub.com/m/v2/packages/const"
	"resturants-hub.com/m/v2/packages/types"
)

// MARK: SsoUserInfo
type SsoUserInfo struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
	Email      string `json:"email"`
}

type User struct {
	Id        int64          `json:"id" db:"id" goqu:"skipinsert"`
	Email     string         `json:"email" db:"email"`
	Role      consts.Role    `json:"role" goqu:"skipinsert"`
	FirstName string         `json:"firstName" db:"first_name"`
	LastName  string         `json:"lastName" db:"last_name"`
	AvatarURL string         `json:"avatarUrl" db:"avatar_url"`
	CreatedAt time.Time      `json:"createdAt" db:"created_at" goqu:"skipinsert"`
	UpdatedAt time.Time      `json:"updatedAt" db:"updated_at" goqu:"skipinsert"`
	DeletedAt types.NullTime `json:"deletedAt" db:"deleted_at" goqu:"skipinsert"`
}

type CreateUserPayload struct {
	ID                   int64  `json:"id" goqu:"skipinsert"`
	Username             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation" goqu:"skipinsert"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PublicListItem struct {
	Id        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	AvatarURL string `json:"avatarUrl"`
}
type AdminListItem struct {
	PublicListItem
	Email string `json:"email"`
}

type AdminDetailItem struct {
	AdminListItem
	Role      consts.Role    `json:"role"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt types.NullTime `json:"deletedAt"`
}
type OwnerDetailItem struct {
	PublicListItem
	Email string      `json:"email"`
	Role  consts.Role `json:"role"`
}

const (
	AdminList    ResponsePayloadType = 0
	AdminDetails                     = 1
	OwnerDetails                     = 2
	PublicList                       = 3
)

type Users []User
type ResponsePayloadType int64

func (user *User) IsAdmin() bool {
	return user.Role == consts.Admin
}

func (user *User) IsManager() bool {
	return user.Role == consts.Manager
}
