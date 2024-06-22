package pages

import (
	"database/sql"
	"time"

	consts "resturants-hub.com/m/v2/packages/const"
	"resturants-hub.com/m/v2/packages/types"
)

// DB representation of the page table
type Page struct {
	Id           int64         `json:"id" db:"id" goqu:"skipinsert,skipupdate"`
	Title        string        `json:"title" db:"title" goqu:"omitempty" validate:"required,min=3,max=50"`
	Slug         string        `json:"slug" db:"slug" validate:"required,min=3,max=50"`
	Excerpt      string        `json:"excerpt" db:"excerpt" goqu:"omitempty" validate:"min=10,max=2000"`
	Body         string        `json:"body" db:"body" goqu:"omitempty" validate:"required,min=100,"`
	Visibility   string        `json:"visibility" db:"visibility" goqu:"omitempty"`
	AuthorId     int64         `json:"authorId" db:"author_id" goqu:"omitempty" validate:"required"`
	RestaurantId types.NullInt `json:"restaurantId" db:"restaurant_id" goqu:"omitempty" validate:"required"`
	ParentPageId types.NullInt `json:"parentPageId" db:"parent_page_id" goqu:"omitempty"`
	CreatedAt    time.Time     `json:"createdAt" db:"created_at" goqu:"skipinsert,skipupdate,omitempty"`
	UpdatedAt    time.Time     `json:"updatedAt" db:"updated_at" goqu:"skipinsert,skipupdate,omitempty"`
	DeletedAt    sql.NullTime  `json:"deletedAt" db:"deleted_at" goqu:"skipupdate,omitempty"`
}

/* Struct for creating new Page */
type CreatePagePayload struct {
	Title        string        `json:"title" db:"title" goqu:"omitempty" validate:"required,min=3,max=50"`
	Slug         string        `json:"slug" db:"slug" validate:"required,min=3,max=50"`
	Excerpt      string        `json:"excerpt" db:"excerpt" goqu:"omitempty" validate:"min=10,max=2000"`
	Body         string        `json:"body" db:"body" goqu:"omitempty" validate:"required,min=100"`
	Visibility   string        `json:"visibility" db:"visibility" goqu:"omitempty"`
	AuthorId     int64         `json:"authorId" db:"author_id" goqu:"omitempty" validate:"required"`
	RestaurantId types.NullInt `json:"restaurantId" db:"restaurant_id" goqu:"omitempty" validate:"required"`
	ParentPageId types.NullInt `json:"parentPageId" db:"parent_page_id" goqu:"omitempty"`
	DeletedAt    sql.NullTime  `json:"deletedAt" db:"deleted_at" goqu:"skipupdate,omitempty"`
}

type OwnerListItem struct {
	Title        string        `json:"title" db:"title"`
	Slug         string        `json:"slug" db:"slug"`
	Excerpt      string        `json:"excerpt" db:"excerpt"`
	Visibility   string        `json:"visibility" db:"visibility"`
	AuthorId     int64         `json:"authorId" db:"author_id"`
	ParentPageId types.NullInt `json:"parentPageId" db:"parent_page_id"`
}

type OwnerDetailItem struct {
	OwnerListItem
	Body string `json:"body" db:"body"`
}
type AdminListItem struct {
	OwnerListItem
	RestaurantId types.NullInt `json:"restaurantId" db:"restaurant_id"`
	DeletedAt    sql.NullTime  `json:"deletedAt" db:"deleted_at"`
}

type AdminDetailItem struct {
	AdminListItem
	Body string `json:"body" db:"body"`
}

type PayloadTypes interface {
	AdminListItem | OwnerListItem | AdminDetailItem | OwnerDetailItem
}

func (page *Page) UpdableAttributes(role consts.Role) []string {
	switch role {
	case consts.Admin:
		return []string{"title", "excerpt", "body", "visibility", "authorId", "restaurantId", "parentPageId", "deletedAt"}
	case consts.Manager:
		return []string{"title", "excerpt", "body", "visibility", "parentPageId", "deletedAt"}
	default:
		return []string{}
	}
}

const (
	AdminList    ViewTypes = 0
	OwnerList              = 1
	AdminDetails           = 2
	OwnerDetails           = 3
)

type Pages []Page
type ViewTypes int64
