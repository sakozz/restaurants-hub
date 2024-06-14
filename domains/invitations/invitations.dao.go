package invitations

import (
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
	"resturants-hub.com/m/v2/database"
	consts "resturants-hub.com/m/v2/packages/const"
	"resturants-hub.com/m/v2/packages/structs"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

// MARK: InvitationsDao
type InvitationsDao interface {
	Create(*CreateInvitationPayload) (*Invitation, rest_errors.RestErr)
	Update(*Invitation, interface{}) (*Invitation, rest_errors.RestErr)
	Get(id *int64) (*Invitation, rest_errors.RestErr)
	Where(params map[string]interface{}) *Invitation
	AuthorizedCollection(url.Values, *structs.BaseUser) (Invitations, rest_errors.RestErr)
}

type connection struct {
	db         *sqlx.DB
	sqlBuilder database.SqlBuilder
}

func NewInvitationDao() InvitationsDao {
	return &connection{
		db:         database.DB,
		sqlBuilder: database.NewSqlBuilder(),
	}
}

func (invitation *Invitation) UpdableAttributes() []string {
	return []string{"expires_at", "role"}
}

func (connection *connection) Create(payload *CreateInvitationPayload) (*Invitation, rest_errors.RestErr) {
	invitation := &Invitation{}
	sqlQuery := connection.sqlBuilder.Insert("invitations", payload)

	row := connection.db.QueryRowx(sqlQuery)
	if row.Err() != nil {
		if uniquenessViolation, constraintName := database.HasUniquenessViolation(row.Err()); uniquenessViolation {
			return nil, rest_errors.InvalidError(ErrorMessage(constraintName))
		}
		return nil, rest_errors.NewInternalServerError(row.Err())
	}

	row.StructScan(invitation)
	return invitation, nil
}

func (connection *connection) Update(invitation *Invitation, payload interface{}) (*Invitation, rest_errors.RestErr) {
	sqlQuery := connection.sqlBuilder.Update("invitations", &invitation.Id, payload)
	row := connection.db.QueryRowx(sqlQuery)
	if row.Err() != nil {
		if uniquenessViolation, constraintName := database.HasUniquenessViolation(row.Err()); uniquenessViolation {
			return nil, rest_errors.InvalidError(ErrorMessage(constraintName))
		}
		return nil, rest_errors.NewInternalServerError(row.Err())
	}
	row.StructScan(invitation)
	return invitation, nil
}

func (connection *connection) Get(id *int64) (*Invitation, rest_errors.RestErr) {
	invitation := &Invitation{}
	query := connection.sqlBuilder.Find("invitations", map[string]interface{}{"id": id})
	err := connection.db.Get(invitation, query)

	if err != nil {
		message := fmt.Sprintf("Sorry, invitation with id %v doesn't exist", *id)
		return nil, rest_errors.NewNotFoundError(message)
	}

	return invitation, nil
}
func (connection *connection) Where(params map[string]interface{}) *Invitation {
	invitation := &Invitation{}

	query := connection.sqlBuilder.SearchBy("invitations", params)
	err := connection.db.Get(invitation, query)
	if err != nil {
		fmt.Println("Error Occured:", err)
		return nil
	}

	return invitation
}

func (connection *connection) AuthorizedCollection(params url.Values, user *structs.BaseUser) (Invitations, rest_errors.RestErr) {
	switch user.Role {
	case consts.Admin:
		return connection.search(params)
	default:
		return Invitations{}, nil
	}
}

func (connection *connection) search(params url.Values) (Invitations, rest_errors.RestErr) {
	var invitations Invitations
	sqlQuery := connection.sqlBuilder.Filter("invitations", params)
	err := connection.db.Select(&invitations, sqlQuery)
	if err != nil {
		return nil, rest_errors.NewNotFoundError(err.Error())
	}

	return invitations, nil
}

func ErrorMessage(errorKey string) string {
	errors := map[string]string{
		"users_username_key": "Username must be unique",
		"users_email_key":    "Email must be unique",
		"user_not_found":     "User is not found",
	}

	return errors[errorKey]
}
