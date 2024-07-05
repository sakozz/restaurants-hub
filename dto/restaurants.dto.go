package dto

import (
	"database/sql"
	"encoding/json"
	"time"

	consts "resturants-hub.com/m/v2/packages/const"
	"resturants-hub.com/m/v2/packages/types"
	"resturants-hub.com/m/v2/serializers"
)

// DB representation of the restaurant table
type Restaurant struct {
	Id            int64                  `json:"id" db:"id" goqu:"skipinsert,skipupdate"`
	ManagerId     int64                  `json:"managerId" db:"manager_id" goqu:"omitempty" validate:"required"`
	Name          string                 `json:"name" db:"name" goqu:"omitempty" validate:"required,min=3,max=50"`
	Description   string                 `json:"description" db:"description" goqu:"omitempty" validate:"required,min=10"`
	Address       types.JsonMap[Address] `json:"address" db:"address" goqu:"omitempty" validate:"required"`
	Email         string                 `json:"email" db:"email" goqu:"omitempty" validate:"required,email"`
	Phone         string                 `json:"phone" db:"phone" goqu:"omitempty" validate:"required"`
	Mobile        string                 `json:"mobile" db:"mobile" goqu:"omitempty"`
	Website       string                 `json:"website" db:"website" goqu:"omitempty"`
	FacebookLink  string                 `json:"facebookLink" db:"facebook_link" goqu:"omitempty"`
	InstagramLink string                 `json:"instagramLink" db:"instagram_link" goqu:"omitempty"`
	CreatedAt     time.Time              `json:"createdAt" db:"created_at" goqu:"skipinsert,skipupdate,omitempty"`
	UpdatedAt     time.Time              `json:"updatedAt" db:"updated_at" goqu:"skipinsert,skipupdate,omitempty"`
	DeletedAt     sql.NullTime           `json:"deletedAt" db:"deleted_at" goqu:"skipupdate,omitempty"`
}
type Address struct {
	Street     string `json:"street" db:"street" validate:"required"`
	City       string `json:"city" db:"city" validate:"required"`
	PostalCode string `json:"postalCode" db:"postal_code" validate:"required"`
	Country    string `json:"country" db:"country" validate:"required"`
}

/* Struct for creating new restaurant */
type CreateRestaurantPayload struct {
	ManagerId     int64                  `json:"managerId" db:"manager_id" validate:"required"`
	Name          string                 `json:"name" db:"name" validate:"required,min=3,max=50"`
	Address       types.JsonMap[Address] `json:"address" db:"address" validate:"required"`
	Description   string                 `json:"description" db:"description" validate:"required,min=10"`
	Email         string                 `json:"email" db:"email" validate:"required,email"`
	Phone         string                 `json:"phone" db:"phone"`
	Mobile        string                 `json:"mobile" db:"mobile" goqu:"omitempty"`
	Website       string                 `json:"website" db:"website" goqu:"omitempty"`
	FacebookLink  string                 `json:"facebookLink" db:"facebook_link" goqu:"omitempty"`
	InstagramLink string                 `json:"instagramLink" db:"instagram_link" goqu:"omitempty"`
}

type AdminRestaurantListItem struct {
	Id          int64                  `json:"id"`
	ManagerId   int64                  `json:"managerId" db:"manager_id"`
	Name        string                 `json:"name" db:"name"`
	Description string                 `json:"description" db:"description"`
	Address     types.JsonMap[Address] `json:"address" db:"address"`
	Email       string                 `json:"email" db:"email"`
	Phone       string                 `json:"phone" db:"phone"`
	CreatedAt   time.Time              `json:"createdAt" db:"created_at"`
}

type AdminRestaurantDetailItem struct {
	Id            int64                  `json:"id"`
	ManagerId     int64                  `json:"managerId" db:"manager_id"`
	Name          string                 `json:"name" db:"name"`
	Description   string                 `json:"description" db:"description"`
	Address       types.JsonMap[Address] `json:"address" db:"address"`
	Email         string                 `json:"email" db:"email"`
	Phone         string                 `json:"phone" db:"phone"`
	Mobile        string                 `json:"mobile" db:"mobile"`
	Website       string                 `json:"website" db:"website"`
	FacebookLink  string                 `json:"facebookLink" db:"facebook_link"`
	InstagramLink string                 `json:"instagramLink" db:"instagram_link"`
	CreatedAt     time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time              `json:"updatedAt" db:"updated_at"`
	DeletedAt     sql.NullTime           `json:"deletedAt" db:"deleted_at"`
}

type OwnerRestaurantListItem struct {
	Id          int64     `json:"id"`
	ManagerId   int64     `json:"managerId" db:"manager_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Address     string    `json:"address" db:"address"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type OwnerRestaurantDetailItem struct {
	Id            int64     `json:"id"`
	ManagerId     int64     `json:"managerId" db:"manager_id"`
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

type RestaurantPayloadTypes interface {
	AdminRestaurantListItem | OwnerRestaurantListItem | AdminRestaurantDetailItem | OwnerRestaurantDetailItem
}

func (restaurant *Restaurant) AdminUpdableAttributes() []string {
	return []string{"managerId", "name", "description", "address", "email", "phone", "mobile", "website", "facebookLink", "instagramLink", "deletedAt"}
}

const (
	AdminList    ViewTypes = 0
	OwnerList              = 1
	AdminDetails           = 2
	OwnerDetails           = 3
)

type Restaurants []Restaurant
type ViewTypes int64

func (record *Restaurant) MemberFor(role consts.Role) interface{} {
	payload, _ := json.Marshal(record)
	switch role {
	case consts.Admin:
		var details AdminRestaurantDetailItem
		json.Unmarshal(payload, &details)
		return serializers.MemberPayload[AdminRestaurantDetailItem]{Id: record.Id, Type: "restaurants", Attributes: details}
	case consts.Manager:
		var details OwnerRestaurantDetailItem
		json.Unmarshal(payload, &details)
		return serializers.MemberPayload[OwnerRestaurantDetailItem]{Id: record.Id, Type: "restaurants", Attributes: details}
	default:
		var listItem OwnerRestaurantListItem
		json.Unmarshal(payload, &listItem)
		return serializers.MemberPayload[OwnerRestaurantListItem]{Id: record.Id, Type: "restaurants", Attributes: listItem}
	}
}

func (restaurants Restaurants) CollectionFor(role consts.Role) []interface{} {
	result := make([]interface{}, len(restaurants))
	for index, record := range restaurants {
		payload, _ := json.Marshal(record)
		switch role {
		case consts.Admin:
			var adminListItem AdminRestaurantListItem
			json.Unmarshal(payload, &adminListItem)
			result[index] = serializers.MemberPayload[AdminRestaurantListItem]{Id: record.Id, Type: "restaurants", Attributes: adminListItem}
		case consts.Manager:
			var details OwnerRestaurantDetailItem
			json.Unmarshal(payload, &details)
			result[index] = serializers.MemberPayload[OwnerRestaurantDetailItem]{Id: record.Id, Type: "restaurants", Attributes: details}
		}
	}
	return result
}
