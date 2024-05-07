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
			return nil, rest_errors.InvalidError(ErrorMessage(constraintName))
		}
		return nil, rest_errors.NewInternalServerError("Server Error", row.Err())
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

func ErrorMessage(errorKey string) string {
	errors := map[string]string{
		"restaurants_name_key":  "name must be unique",
		"restaurants_email_key": "Email must be unique",
		"restaurant_not_found":  "Restaurant is not found",
	}

	return errors[errorKey]
}
