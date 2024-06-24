package users

import (
	"encoding/json"

	"resturants-hub.com/m/v2/jsonapi"
	consts "resturants-hub.com/m/v2/packages/const"
)

func (user *User) MemberFor(role consts.Role) interface{} {
	payload, _ := json.Marshal(user)
	switch role {

	case consts.Admin:
		var details AdminDetailItem
		json.Unmarshal(payload, &details)
		return jsonapi.MemberPayload[AdminDetailItem]{Id: user.Id, Type: "users", Attributes: details}
	default:
		var details OwnerDetailItem
		json.Unmarshal(payload, &details)
		return jsonapi.MemberPayload[OwnerDetailItem]{Id: user.Id, Type: "users", Attributes: details}
	}
}

func (users Users) CollectionFor(role consts.Role) []interface{} {
	result := make([]interface{}, len(users))
	for index, record := range users {
		payload, _ := json.Marshal(record)
		switch role {
		case consts.Admin:
			var adminListItem AdminListItem
			json.Unmarshal(payload, &adminListItem)
			result[index] = jsonapi.MemberPayload[AdminListItem]{Id: record.Id, Type: "users", Attributes: adminListItem}
		default:
			var publicListItem PublicListItem
			json.Unmarshal(payload, &publicListItem)
			result[index] = jsonapi.MemberPayload[PublicListItem]{Id: record.Id, Type: "users", Attributes: publicListItem}
		}
	}
	return result
}
