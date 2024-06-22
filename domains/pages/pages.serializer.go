package pages

import (
	"encoding/json"

	"resturants-hub.com/m/v2/jsonapi"
	consts "resturants-hub.com/m/v2/packages/const"
)

func (record *Page) MemberFor(role consts.Role) interface{} {
	payload, _ := json.Marshal(record)
	switch role {
	case consts.Admin:
		var details AdminDetailItem
		json.Unmarshal(payload, &details)
		return jsonapi.MemberPayload[AdminDetailItem]{Id: record.Id, Type: "pages", Attributes: details}
	case consts.Manager:
		var details OwnerDetailItem
		json.Unmarshal(payload, &details)
		return jsonapi.MemberPayload[OwnerDetailItem]{Id: record.Id, Type: "pages", Attributes: details}
	default:
		var details PublicItem
		json.Unmarshal(payload, &details)
		return jsonapi.MemberPayload[PublicItem]{Id: record.Id, Type: "pages", Attributes: details}
	}
}

func (pages Pages) CollectionFor(role consts.Role) []interface{} {
	result := make([]interface{}, len(pages))
	for index, record := range pages {
		payload, _ := json.Marshal(record)
		switch role {
		case consts.Admin:
			var item AdminListItem
			json.Unmarshal(payload, &item)
			result[index] = jsonapi.MemberPayload[AdminListItem]{Id: record.Id, Type: "pages", Attributes: item}
		case consts.Manager:
			var item OwnerListItem
			json.Unmarshal(payload, &item)
			result[index] = jsonapi.MemberPayload[OwnerListItem]{Id: record.Id, Type: "pages", Attributes: item}
		default:
			var item PublicItem
			json.Unmarshal(payload, &item)
			result[index] = jsonapi.MemberPayload[PublicItem]{Id: record.Id, Type: "pages", Attributes: item}
		}
	}
	return result
}
