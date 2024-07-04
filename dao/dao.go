package dao

import (
	"github.com/jmoiron/sqlx"
	"resturants-hub.com/m/v2/database"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

type connection struct {
	db         *sqlx.DB
	sqlBuilder database.SqlBuilder
}

func ErrorMessage(errorKey string) string {
	errors := map[string]string{
		"users_username_key": "Username must be unique",
		"users_email_key":    "Email must be unique",
		"user_not_found":     "User is not found",
	}

	return errors[errorKey]
}

func UniquenessErrors(errorKey string) *rest_errors.ValidationErrs {
	causes := rest_errors.ValidationErrs{}
	errKeyMaps := map[string]string{
		"restaurants_name_key":  "name",
		"restaurants_email_key": "email",
		"restaurants_phone_key": "phone",
		"fk_user":               "userId",
		"pages_name_key":        "name",
		"pages_email_key":       "email",
		"pages_phone_key":       "phone",
	}

	attr := errKeyMaps[errorKey]
	causes[attr] = []interface{}{
		rest_errors.FormattedDbValidationError(attr, "uniqueness"),
	}
	return &causes
}
