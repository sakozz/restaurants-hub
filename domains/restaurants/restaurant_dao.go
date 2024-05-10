package restaurants

import (
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
	"resturants-hub.com/m/v2/database"
	rest_errors "resturants-hub.com/m/v2/utils"
)

type connection struct {
	db         *sqlx.DB
	sqlBuilder database.SqlBuilder
}

type RestaurantDao interface {
	Create(*CreateRestaurantPayload) (*Restaurant, rest_errors.RestErr)
	Search(url.Values) (Restaurants, rest_errors.RestErr)
}

func NewRestaurantDao() RestaurantDao {
	return &connection{
		db:         database.DB,
		sqlBuilder: database.NewSqlBuilder(),
	}
}

func (connection *connection) Create(payload *CreateRestaurantPayload) (*Restaurant, rest_errors.RestErr) {
	restaurant := &Restaurant{}
	sqlQuery := connection.sqlBuilder.Insert("restaurants", payload)
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

func (connection *connection) Search(params url.Values) (Restaurants, rest_errors.RestErr) {
	var restaurants Restaurants
	sqlQuery := connection.sqlBuilder.Filter("restaurants", params)
	err := connection.db.Select(&restaurants, sqlQuery)
	if err != nil {
		return nil, rest_errors.NewNotFoundError(err.Error())
	}

	return restaurants, nil
}

func UniquenessErrors(errorKey string) *rest_errors.ValidationErrs {
	causes := rest_errors.ValidationErrs{}
	errKeyMaps := map[string]string{
		"restaurants_name_key":  "name",
		"restaurants_email_key": "email",
		"restaurants_phone_key": "phone",
	}

	attr := errKeyMaps[errorKey]
	causes[attr] = []interface{}{
		rest_errors.FormattedDbValidationError(attr, "uniqueness"),
	}
	return &causes
}
