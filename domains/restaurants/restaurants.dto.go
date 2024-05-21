package restaurants

import (
	"database/sql"
	"time"

	data "resturants-hub.com/m/v2/data_types"
)

// DB representation of the restaurant table
type Restaurant struct {
	Id            int64                 `json:"id" db:"id" goqu:"skipinsert,skipupdate"`
	ProfileId     int64                 `json:"profileId" db:"profile_id" goqu:"omitempty" validate:"required"`
	Name          string                `json:"name" db:"name" goqu:"omitempty" validate:"required,min=3,max=50"`
	Description   string                `json:"description" db:"description" goqu:"omitempty" validate:"required,min=10"`
	Address       data.JsonMap[Address] `json:"address" db:"address" goqu:"omitempty" validate:"required"`
	Email         string                `json:"email" db:"email" goqu:"omitempty" validate:"required,email"`
	Phone         string                `json:"phone" db:"phone" goqu:"omitempty" validate:"required"`
	Mobile        string                `json:"mobile" db:"mobile" goqu:"omitempty"`
	Website       string                `json:"website" db:"website" goqu:"omitempty"`
	FacebookLink  string                `json:"facebookLink" db:"facebook_link" goqu:"omitempty"`
	InstagramLink string                `json:"instagramLink" db:"instagram_link" goqu:"omitempty"`
	CreatedAt     time.Time             `json:"createdAt" db:"created_at" goqu:"skipinsert,skipupdate,omitempty"`
	UpdatedAt     time.Time             `json:"updatedAt" db:"updated_at" goqu:"skipinsert,skipupdate,omitempty"`
	DeletedAt     sql.NullTime          `json:"deletedAt" db:"deleted_at" goqu:"skipupdate,omitempty"`
}
type Address struct {
	Street     string `json:"street" db:"street" validate:"required"`
	City       string `json:"city" db:"city" validate:"required"`
	PostalCode string `json:"code" db:"postal_code" validate:"required"`
	Country    string `json:"country" db:"country" validate:"required"`
}

/* Struct for creating new restaurant */
type CreateRestaurantPayload struct {
	ProfileId     int64                 `json:"profileId" db:"profile_id" validate:"required"`
	Name          string                `json:"name" db:"name" validate:"required,min=3,max=50"`
	Address       data.JsonMap[Address] `json:"address" db:"address" validate:"required"`
	Description   string                `json:"description" db:"description" validate:"required,min=10"`
	Email         string                `json:"email" db:"email" validate:"required,email"`
	Phone         string                `json:"phone" db:"phone"`
	Mobile        string                `json:"mobile" db:"mobile" goqu:"omitempty"`
	Website       string                `json:"website" db:"website" goqu:"omitempty"`
	FacebookLink  string                `json:"facebookLink" db:"facebook_link" goqu:"omitempty"`
	InstagramLink string                `json:"instagramLink" db:"instagram_link" goqu:"omitempty"`
}

type AdminListItem struct {
	Id          int64                 `json:"id"`
	ProfileId   int64                 `json:"profileId" db:"profile_id"`
	Name        string                `json:"name" db:"name"`
	Description string                `json:"description" db:"description"`
	Address     data.JsonMap[Address] `json:"address" db:"address"`
	Email       string                `json:"email" db:"email"`
	Phone       string                `json:"phone" db:"phone"`
	CreatedAt   time.Time             `json:"createdAt" db:"created_at"`
}

type AdminDetailItem struct {
	Id            int64                 `json:"id"`
	ProfileId     int64                 `json:"profileId" db:"profile_id"`
	Name          string                `json:"name" db:"name"`
	Description   string                `json:"description" db:"description"`
	Address       data.JsonMap[Address] `json:"address" db:"address"`
	Email         string                `json:"email" db:"email"`
	Phone         string                `json:"phone" db:"phone"`
	Mobile        string                `json:"mobile" db:"mobile"`
	Website       string                `json:"website" db:"website"`
	FacebookLink  string                `json:"facebookLink" db:"facebook_link"`
	InstagramLink string                `json:"instagramLink" db:"instagram_link"`
	CreatedAt     time.Time             `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time             `json:"updatedAt" db:"updated_at"`
	DeletedAt     sql.NullTime          `json:"deletedAt" db:"deleted_at"`
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

type PayloadTypes interface {
	AdminListItem | OwnerListItem | AdminDetailItem | OwnerDetailItem
}

func (restaurant *Restaurant) AdminUpdableAttributes() []string {
	return []string{"profileId", "name", "description", "address", "email", "phone", "mobile", "website", "facebookLink", "instagramLink", "deletedAt"}
}

const (
	AdminList    ViewTypes = 0
	OwnerList              = 1
	AdminDetails           = 2
	OwnerDetails           = 3
)

type Restaurants []Restaurant
type ViewTypes int64
