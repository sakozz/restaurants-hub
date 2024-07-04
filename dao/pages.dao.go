package dao

import (
	"fmt"
	"net/url"

	"github.com/gosimple/slug"
	"github.com/mitchellh/mapstructure"
	"resturants-hub.com/m/v2/database"
	"resturants-hub.com/m/v2/dto"
	consts "resturants-hub.com/m/v2/packages/const"
	"resturants-hub.com/m/v2/packages/structs"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type PagesDao interface {
	Create(*dto.CreatePagePayload) (*dto.Page, rest_errors.RestErr)
	Search(url.Values) (dto.Pages, rest_errors.RestErr)
	AuthorizedCollection(url.Values, *structs.BaseUser) (dto.Pages, rest_errors.RestErr)
	Get(slug *string) (*dto.Page, rest_errors.RestErr)
	Update(*dto.Page, interface{}) (*dto.Page, rest_errors.RestErr)
	GenerateSlug(string) string
}

func NewPageDao() PagesDao {
	return &connection{
		db:         database.DB,
		sqlBuilder: database.NewSqlBuilder(),
	}
}

func (connection *connection) Create(payload *dto.CreatePagePayload) (*dto.Page, rest_errors.RestErr) {
	restaurant := &dto.Page{}
	sqlQuery := connection.sqlBuilder.Insert("pages", payload)
	row := connection.db.QueryRowx(sqlQuery)
	if row.Err() != nil {
		fmt.Println(row.Err())
		if uniquenessViolation, constraintName := database.HasUniquenessViolation(row.Err()); uniquenessViolation {
			return nil, rest_errors.NewValidationError(UniquenessErrors(constraintName))
		}
		return nil, rest_errors.NewInternalServerError(row.Err())
	}

	row.StructScan(restaurant)
	return restaurant, nil
}

func (connection *connection) Get(slug *string) (*dto.Page, rest_errors.RestErr) {
	restaurant := &dto.Page{}
	query := connection.sqlBuilder.Find("pages", map[string]interface{}{"slug": slug})
	err := connection.db.Get(restaurant, query)

	if err != nil {
		message := fmt.Sprintf("Sorry, the record with slug %v doesn't exist", *slug)
		return nil, rest_errors.NewNotFoundError(message)
	}

	return restaurant, nil
}

func (connection *connection) Search(params url.Values) (dto.Pages, rest_errors.RestErr) {
	var pages dto.Pages
	sqlQuery := connection.sqlBuilder.Filter("pages", params)
	err := connection.db.Select(&pages, sqlQuery)
	if err != nil {
		return nil, rest_errors.NewNotFoundError(err.Error())
	}

	return pages, nil
}

func (connection *connection) GenerateSlug(title string) string {
	slug := slug.Make(title)
	slugExists := true
	for slugExists == true {
		_, err := connection.Get(&slug)
		/* If err == nil (the record with the generated slug exist) */
		slugExists = err == nil
		if slugExists {
			slug = slug + "-1"
		}
	}
	return slug
}

func (connection *connection) AuthorizedCollection(params url.Values, user *structs.BaseUser) (dto.Pages, rest_errors.RestErr) {
	switch user.Role {
	case consts.Admin:
		return connection.Search(params)
	case consts.Manager:
		params.Add("author_id", fmt.Sprint(user.Id))
		return connection.Search(params)
	default:
		return dto.Pages{}, nil
	}
}

func (connection *connection) Update(page *dto.Page, payload interface{}) (*dto.Page, rest_errors.RestErr) {
	// Convert payload to Page struct: this is to ensure that attribute names are mapped with db column names
	payloadPage := &dto.Page{}
	mapstructure.Decode(payload, payloadPage)

	sqlQuery := connection.sqlBuilder.Update("pages", &page.Id, payloadPage)
	row := connection.db.QueryRowx(sqlQuery)
	if row.Err() != nil {
		if uniquenessViolation, constraintName := database.HasUniquenessViolation(row.Err()); uniquenessViolation {
			return nil, rest_errors.NewValidationError(UniquenessErrors(constraintName))
		}
		return nil, rest_errors.NewInternalServerError(row.Err())
	}
	row.StructScan(page)
	return page, nil
}
