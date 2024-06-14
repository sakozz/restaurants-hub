package invitations

import (
	"time"

	consts "resturants-hub.com/m/v2/packages/const"
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
