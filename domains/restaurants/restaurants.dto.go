package restaurants

import (
	"database/sql"
	"time"

	rest_errors "resturants-hub.com/m/v2/utils"
)

// DB representation of the restaurant table
type Restaurant struct {
	Id            int64        `json:"id"`
	ProfileId     int64        `json:"profileId" db:"profile_id"`
	Name          string       `json:"name" db:"name"`
	Description   string       `json:"description" db:"description"`
	Address       string       `json:"address" db:"address"`
	Email         string       `json:"email" db:"email"`
	Phone         string       `json:"phone" db:"phone"`
	Mobile        string       `json:"mobile" db:"mobile"`
	Website       string       `json:"website" db:"website"`
	FacebookLink  string       `json:"facebookLink" db:"facebook_link"`
	InstagramLink string       `json:"instagramLink" db:"instagram_link"`
	CreatedAt     time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time    `json:"updatedAt" db:"updated_at"`
	DeletedAt     sql.NullTime `json:"deletedAt" db:"deleted_at"`
}

/* Struct for creating new restaurant */
type CreateRestaurantPayload struct {
	ProfileId     int64  `json:"profileId" db:"profile_id" validate:"required"`
	Name          string `json:"name" db:"name" validate:"required,min=3,max=50"`
	Description   string `json:"description" db:"description" validate:"required,min=10"`
	Email         string `json:"email" db:"email" validate:"required,email"`
	Phone         string `json:"phone" db:"phone"`
	Mobile        string `json:"mobile" db:"mobile" goqu:"omitempty"`
	Website       string `json:"website" db:"website" goqu:"omitempty"`
	FacebookLink  string `json:"facebookLink" db:"facebook_ink" goqu:"omitempty"`
	InstagramLink string `json:"instagramLink" db:"instagram_link" goqu:"omitempty"`
}

type AdminListItem struct {
	Id          int64     `json:"id"`
	ProfileId   int64     `json:"profileId" db:"profile_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Address     string    `json:"address" db:"address"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type AdminDetailItem struct {
	Id            int64        `json:"id"`
	ProfileId     int64        `json:"profileId" db:"profile_id"`
	Name          string       `json:"name" db:"name"`
	Description   string       `json:"description" db:"description"`
	Address       string       `json:"address" db:"address"`
	Email         string       `json:"email" db:"email"`
	Phone         string       `json:"phone" db:"phone"`
	Mobile        string       `json:"mobile" db:"mobile"`
	Website       string       `json:"website" db:"website"`
	FacebookLink  string       `json:"facebookLink" db:"facebook_link"`
	InstagramLink string       `json:"instagramLink" db:"instagram_link"`
	CreatedAt     time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time    `json:"updatedAt" db:"updated_at"`
	DeletedAt     sql.NullTime `json:"deletedAt" db:"deleted_at"`
}

type OwnerListItem struct {
	Id          int64     `json:"id"`
	ProfileId   int64     `json:"profileId" db:"profile_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Address     string    `json:"address" db:"address"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type OwnerDetailItem struct {
	Id            int64     `json:"id"`
	ProfileId     int64     `json:"profileId" db:"profile_id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	Address       string    `json:"address" db:"address"`
	Email         string    `json:"email" db:"email"`
	Phone         string    `json:"phone" db:"phone"`
	Mobile        string    `json:"mobile" db:"mobile"`
	Website       string    `json:"website" db:"website"`
	FacebookLink  string    `json:"facebookLink" db:"facebook_link"`
	InstagramLink string    `json:"instagramLink" db:"instagram_link"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

func (newRestaurant *CreateRestaurantPayload) Validate() rest_errors.RestErr {
	return nil
}

type PayloadTypes interface {
	AdminListItem | OwnerListItem | AdminDetailItem | OwnerDetailItem
}
type Payload[T PayloadTypes] struct {
	Id         int64  `json:"id"`
	Type       string `json:"type"`
	Attributes T      `json:"attributes"`
}

const (
	AdminList    ViewTypes = 0
	OwnerList              = 1
	AdminDetails           = 2
	OwnerDetails           = 3
)

type Restaurants []Restaurant
type ViewTypes int64
