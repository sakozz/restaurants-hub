package restaurants

import (
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
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
	Get(id *int64) (*Restaurant, rest_errors.RestErr)
	Update(*Restaurant, interface{}) (*Restaurant, rest_errors.RestErr)
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
func (connection *connection) Get(id *int64) (*Restaurant, rest_errors.RestErr) {
	restaurant := &Restaurant{}
	query := connection.sqlBuilder.Find("restaurants", map[string]interface{}{"id": id})
	err := connection.db.Get(restaurant, query)

	if err != nil {
		message := fmt.Sprintf("Sorry, the record with id %v doesn't exist", *id)
		return nil, rest_errors.NewNotFoundError(message)
	}

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

func (connection *connection) Update(restaurant *Restaurant, payload interface{}) (*Restaurant, rest_errors.RestErr) {
	// Convert payload to Restaurant struct: this is to ensure that attribute names are mapped with db column names
	payloadRestaurant := &Restaurant{}
	mapstructure.Decode(payload, payloadRestaurant)

	sqlQuery := connection.sqlBuilder.Update("restaurants", &restaurant.Id, payloadRestaurant)
	row := connection.db.QueryRowx(sqlQuery)
	if row.Err() != nil {
		if uniquenessViolation, constraintName := database.HasUniquenessViolation(row.Err()); uniquenessViolation {
			return nil, rest_errors.NewValidationError(UniquenessErrors(constraintName))
		}
		return nil, rest_errors.NewInternalServerError(row.Err())
	}
	row.StructScan(restaurant)
	return restaurant, nil
}

func UniquenessErrors(errorKey string) *rest_errors.ValidationErrs {
	causes := rest_errors.ValidationErrs{}
	errKeyMaps := map[string]string{
		"restaurants_name_key":  "name",
		"restaurants_email_key": "email",
		"restaurants_phone_key": "phone",
		"fk_profile":            "profileId",
	}

	attr := errKeyMaps[errorKey]
	causes[attr] = []interface{}{
		rest_errors.FormattedDbValidationError(attr, "uniqueness"),
	}
	return &causes
}
