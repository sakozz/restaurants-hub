package users

import (
	"database/sql"
	"time"
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
	ID        int64        `json:"id" db:"id" goqu:"skipinsert"`
	Email     string       `json:"email" db:"email"`
	Role      string       `json:"role" goqu:"skipinsert"`
	FirstName string       `json:"firstName" db:"first_name"`
	LastName  string       `json:"lastName" db:"last_name"`
	AvatarURL string       `json:"avatarUrl" db:"avatar_url"`
	CreatedAt sql.NullTime `json:"createdAt" db:"created_at" goqu:"skipinsert"`
	UpdatedAt time.Time    `json:"updatedAt" db:"updated_at" goqu:"skipinsert"`
	DeletedAt sql.NullTime `json:"deletedAt" db:"deleted_at" goqu:"skipinsert"`
}

type SearchUserParams struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   string `json:"status"`
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

type PublicUser struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	AvatarURL string `json:"avatarUrl"`
}

type PrivateUser struct {
	ID        int64        `json:"id"`
	Email     string       `json:"email"`
	FirstName string       `json:"firstName"`
	LastName  string       `json:"lastName"`
	AvatarURL string       `json:"avatarUrl"`
	CreatedAt sql.NullTime `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	DeletedAt sql.NullTime `json:"deletedAt"`
}

// MARK: PayloadTypes
type PayloadTypes interface {
	PublicUser | PrivateUser | SearchUserParams
}

/* Data format for returning payload */
// MARK: UserPayload
type UserPayload[T PayloadTypes] struct {
	Id         int64  `json:"id"`
	Type       string `json:"type"`
	Attributes T      `json:"attributes"`
}

type Users []User
type AuthType int64
