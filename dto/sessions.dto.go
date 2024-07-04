package dto

import "time"

type Session struct {
	Id                int64     `json:"id" db:"id" goqu:"skipinsert"`
	UserId            int64     `json:"userId" db:"user_id" goqu:"omitempty"`
	Provider          string    `json:"provider" db:"provider" goqu:"omitempty"`
	Email             string    `json:"email" db:"email" goqu:"omitempty"`
	AccessToken       string    `json:"accessToken" db:"access_token" goqu:"omitempty"`
	AccessTokenSecret string    `json:"accessTokenSecret" db:"access_token_secret" goqu:"omitempty"`
	RefreshToken      string    `json:"refreshToken" db:"refresh_token" goqu:"omitempty"`
	ExpiresAt         time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at" goqu:"skipinsert omitempty"`
	UpdatedAt         time.Time `json:"updatedAt" db:"updated_at" goqu:"skipinsert"`
	IDToken           string    `json:"idToken" db:"id_token"`
}

type SsoUserInfo struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
	Email      string `json:"email"`
}
