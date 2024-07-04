package dao

import (
	"fmt"
	"net/url"

	"github.com/mitchellh/mapstructure"
	"resturants-hub.com/m/v2/database"
	"resturants-hub.com/m/v2/dto"
	consts "resturants-hub.com/m/v2/packages/const"
	"resturants-hub.com/m/v2/packages/structs"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type RestaurantDao interface {
	CreateRestaurant(*dto.CreateRestaurantPayload) (*dto.Restaurant, rest_errors.RestErr)
	SearchRestaurants(url.Values) (dto.Restaurants, rest_errors.RestErr)
	AuthorizedRestaurantCollection(url.Values, *structs.BaseUser) (dto.Restaurants, rest_errors.RestErr)
	GetRestaurant(id *int64) (*dto.Restaurant, rest_errors.RestErr)
	RestaurantByOwnerId(*int64) (*dto.Restaurant, rest_errors.RestErr)
	UpdateRestaurant(*dto.Restaurant, interface{}) (*dto.Restaurant, rest_errors.RestErr)
}

func NewRestaurantDao() RestaurantDao {
	return &connection{
		db:         database.DB,
		sqlBuilder: database.NewSqlBuilder(),
	}
}

func (connection *connection) CreateRestaurant(payload *dto.CreateRestaurantPayload) (*dto.Restaurant, rest_errors.RestErr) {
	restaurant := &dto.Restaurant{}
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

func (connection *connection) GetRestaurant(id *int64) (*dto.Restaurant, rest_errors.RestErr) {
	restaurant := &dto.Restaurant{}
	query := connection.sqlBuilder.Find("restaurants", map[string]interface{}{"id": id})
	err := connection.db.Get(restaurant, query)

	if err != nil {
		message := fmt.Sprintf("Sorry, the record with id %v doesn't exist", *id)
		return nil, rest_errors.NewNotFoundError(message)
	}

	return restaurant, nil
}

func (connection *connection) RestaurantByOwnerId(id *int64) (*dto.Restaurant, rest_errors.RestErr) {
	restaurant := &dto.Restaurant{}
	params := map[string]interface{}{"manager_id": id}

	query := connection.sqlBuilder.SearchBy(string(consts.Restaurants), params)
	err := connection.db.Get(restaurant, query)
	if err != nil {
		message := fmt.Sprintf("Sorry, the record with id doesn't exist")
		return nil, rest_errors.NewNotFoundError(message)
	}

	return restaurant, nil
}

func (connection *connection) SearchRestaurants(params url.Values) (dto.Restaurants, rest_errors.RestErr) {
	var restaurants dto.Restaurants
	sqlQuery := connection.sqlBuilder.Filter("restaurants", params)
	err := connection.db.Select(&restaurants, sqlQuery)
	if err != nil {
		return nil, rest_errors.NewNotFoundError(err.Error())
	}

	return restaurants, nil
}

func (connection *connection) AuthorizedRestaurantCollection(params url.Values, user *structs.BaseUser) (dto.Restaurants, rest_errors.RestErr) {
	switch user.Role {
	case consts.Admin:
		return connection.SearchRestaurants(params)
	case consts.Manager:
		params.Add("user_id", fmt.Sprint(user.Id))
		return connection.SearchRestaurants(params)
	default:
		return dto.Restaurants{}, nil
	}
}

func (connection *connection) UpdateRestaurant(restaurant *dto.Restaurant, payload interface{}) (*dto.Restaurant, rest_errors.RestErr) {
	// Convert payload to Restaurant struct: this is to ensure that attribute names are mapped with db column names
	payloadRestaurant := &dto.Restaurant{}
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
