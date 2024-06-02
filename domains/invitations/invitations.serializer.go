package invitations

import (
	"encoding/json"

	"resturants-hub.com/m/v2/jsonapi"
)

func (invitation *Invitation) MemberFor() interface{} {
	payload, _ := json.Marshal(invitation)
	var details Invitation
	json.Unmarshal(payload, &details)
	return jsonapi.MemberPayload[Invitation]{Id: invitation.Id, Type: "invitations", Attributes: details}

}

func (invitations Invitations) CollectionFor() []interface{} {
	result := make([]interface{}, len(invitations))
	for index, record := range invitations {
		result[index] = record.MemberFor()
	}
	return result
}
