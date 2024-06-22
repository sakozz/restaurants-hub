package pages

import (
	"fmt"
	"net/url"

	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
	"resturants-hub.com/m/v2/database"
	consts "resturants-hub.com/m/v2/packages/const"
	"resturants-hub.com/m/v2/packages/structs"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type connection struct {
	db         *sqlx.DB
	sqlBuilder database.SqlBuilder
}

type PagesDao interface {
	Create(*CreatePagePayload) (*Page, rest_errors.RestErr)
	Search(url.Values) (Pages, rest_errors.RestErr)
	AuthorizedCollection(url.Values, *structs.BaseUser) (Pages, rest_errors.RestErr)
	Get(slug *string) (*Page, rest_errors.RestErr)
	Update(*Page, interface{}) (*Page, rest_errors.RestErr)
	GenerateSlug(string) string
}

func NewPageDao() PagesDao {
	return &connection{
		db:         database.DB,
		sqlBuilder: database.NewSqlBuilder(),
	}
}

func (connection *connection) Create(payload *CreatePagePayload) (*Page, rest_errors.RestErr) {
	restaurant := &Page{}
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

func (connection *connection) Get(slug *string) (*Page, rest_errors.RestErr) {
	restaurant := &Page{}
	query := connection.sqlBuilder.Find("pages", map[string]interface{}{"slug": slug})
	err := connection.db.Get(restaurant, query)

	if err != nil {
		message := fmt.Sprintf("Sorry, the record with slug %v doesn't exist", *slug)
		return nil, rest_errors.NewNotFoundError(message)
	}

	return restaurant, nil
}

func (connection *connection) Search(params url.Values) (Pages, rest_errors.RestErr) {
	var pages Pages
	sqlQuery := connection.sqlBuilder.Filter("pages", params)
	err := connection.db.Select(&pages, sqlQuery)
	if err != nil {
		return nil, rest_errors.NewNotFoundError(err.Error())
	}

	return pages, nil
}

func (connection *connection) GenerateSlug(title string) string {
	slug:=slug.Make(title)
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

func (connection *connection) AuthorizedCollection(params url.Values, user *structs.BaseUser) (Pages, rest_errors.RestErr) {
	switch user.Role {
	case consts.Admin:
		return connection.Search(params)
	case consts.Manager:
		params.Add("author_id", fmt.Sprint(user.Id))
		return connection.Search(params)
	default:
		return Pages{}, nil
	}
}

func (connection *connection) Update(page *Page, payload interface{}) (*Page, rest_errors.RestErr) {
	// Convert payload to Page struct: this is to ensure that attribute names are mapped with db column names
	payloadPage := &Page{}
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

func UniquenessErrors(errorKey string) *rest_errors.ValidationErrs {
	causes := rest_errors.ValidationErrs{}
	errKeyMaps := map[string]string{
		"pages_name_key":  "name",
		"pages_email_key": "email",
		"pages_phone_key": "phone",
		"fk_user":         "userId",
	}

	attr := errKeyMaps[errorKey]
	causes[attr] = []interface{}{
		rest_errors.FormattedDbValidationError(attr, "uniqueness"),
	}
	return &causes
}
