package dto

import (
	"encoding/json"
	"time"

	consts "resturants-hub.com/m/v2/packages/const"
	"resturants-hub.com/m/v2/serializers"
)

type Invitation struct {
	Id        int64       `json:"id" db:"id" goqu:"skipinsert"`
	Email     string      `json:"email" db:"email"`
	Token     string      `json:"token" db:"token"`
	Role      consts.Role `json:"role" goqu:"skipinsert"`
	CreatedAt time.Time   `json:"createdAt" db:"created_at" goqu:"skipinsert"`
	UpdatedAt time.Time   `json:"updatedAt" db:"updated_at" goqu:"skipinsert"`
	ExpiresAt time.Time   `json:"expiresAt" db:"expires_at" goqu:"skipinsert"`
}

type CreateInvitationPayload struct {
	Email string      `json:"email" db:"email" validate:"required"`
	Token string      `json:"token" db:"token"`
	Role  consts.Role `json:"role" db:"role" validate:"required"`
}

type Invitations []Invitation

func (invitation *Invitation) IsValid() bool {
	return invitation.ExpiresAt.After(time.Now())
}

func (invitation *Invitation) UpdableAttributes() []string {
	return []string{"expires_at", "role"}
}

func (invitation *Invitation) MemberFor() interface{} {
	payload, _ := json.Marshal(invitation)
	var details Invitation
	json.Unmarshal(payload, &details)
	return serializers.MemberPayload[Invitation]{Id: invitation.Id, Type: "invitations", Attributes: details}

}

func (invitations Invitations) CollectionFor() []interface{} {
	result := make([]interface{}, len(invitations))
	for index, record := range invitations {
		result[index] = record.MemberFor()
	}
	return result
}
