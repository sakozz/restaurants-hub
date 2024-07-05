package dao

import (
	"fmt"
	"net/url"

	"resturants-hub.com/m/v2/database"
	"resturants-hub.com/m/v2/dto"
	consts "resturants-hub.com/m/v2/packages/const"
	rest_errors "resturants-hub.com/m/v2/packages/utils"
)

// MARK: InvitationsDao
type InvitationsDao interface {
	CreateInvitation(*dto.CreateInvitationPayload) (*dto.Invitation, rest_errors.RestErr)
	UpdateInvitation(*dto.Invitation, interface{}) (*dto.Invitation, rest_errors.RestErr)
	GetInvitation(id *int64) (*dto.Invitation, rest_errors.RestErr)
	SearchInvitations(params map[string]interface{}) *dto.Invitation
	AuthorizedInvitationsCollection(url.Values, *dto.BaseUser) (dto.Invitations, rest_errors.RestErr)
}

func NewInvitationDao() InvitationsDao {
	return &connection{
		db:         database.DB,
		sqlBuilder: database.NewSqlBuilder(),
	}
}

func (connection *connection) CreateInvitation(payload *dto.CreateInvitationPayload) (*dto.Invitation, rest_errors.RestErr) {
	invitation := &dto.Invitation{}
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

func (connection *connection) UpdateInvitation(invitation *dto.Invitation, payload interface{}) (*dto.Invitation, rest_errors.RestErr) {
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

func (connection *connection) GetInvitation(id *int64) (*dto.Invitation, rest_errors.RestErr) {
	invitation := &dto.Invitation{}
	query := connection.sqlBuilder.Find("invitations", map[string]interface{}{"id": id})
	err := connection.db.Get(invitation, query)

	if err != nil {
		message := fmt.Sprintf("Sorry, invitation with id %v doesn't exist", *id)
		return nil, rest_errors.NewNotFoundError(message)
	}

	return invitation, nil
}
func (connection *connection) SearchInvitations(params map[string]interface{}) *dto.Invitation {
	invitation := &dto.Invitation{}

	query := connection.sqlBuilder.SearchBy("invitations", params)
	err := connection.db.Get(invitation, query)
	if err != nil {
		fmt.Println("Error Occured:", err)
		return nil
	}

	return invitation
}

func (connection *connection) AuthorizedInvitationsCollection(params url.Values, user *dto.BaseUser) (dto.Invitations, rest_errors.RestErr) {
	switch user.Role {
	case consts.Admin:
		return connection.search(params)
	default:
		return dto.Invitations{}, nil
	}
}

func (connection *connection) search(params url.Values) (dto.Invitations, rest_errors.RestErr) {
	var invitations dto.Invitations
	sqlQuery := connection.sqlBuilder.Filter("invitations", params)
	err := connection.db.Select(&invitations, sqlQuery)
	if err != nil {
		return nil, rest_errors.NewNotFoundError(err.Error())
	}

	return invitations, nil
}
