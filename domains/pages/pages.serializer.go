package pages

import (
	"encoding/json"

	"resturants-hub.com/m/v2/jsonapi"
)

func (record *Page) MemberFor(view ViewTypes) interface{} {
	payload, _ := json.Marshal(record)
	switch view {
	case AdminList:
		var adminListItem AdminListItem
		json.Unmarshal(payload, &adminListItem)
		return jsonapi.MemberPayload[AdminListItem]{Id: record.Id, Type: "pages", Attributes: adminListItem}
	case AdminDetails:
		var details AdminDetailItem
		json.Unmarshal(payload, &details)
		return jsonapi.MemberPayload[AdminDetailItem]{Id: record.Id, Type: "pages", Attributes: details}
	default:
		var adminListItem AdminListItem
		json.Unmarshal(payload, &adminListItem)
		return jsonapi.MemberPayload[AdminListItem]{Id: record.Id, Type: "pages", Attributes: adminListItem}
	}
}

func (pages Pages) CollectionFor(view ViewTypes) []interface{} {
	result := make([]interface{}, len(pages))
	for index, record := range pages {
		payload, _ := json.Marshal(record)
		switch view {
		case AdminList:
			var adminListItem AdminListItem
			json.Unmarshal(payload, &adminListItem)
			result[index] = jsonapi.MemberPayload[AdminListItem]{Id: record.Id, Type: "pages", Attributes: adminListItem}
		case AdminDetails:
			var details AdminDetailItem
			json.Unmarshal(payload, &details)
			result[index] = jsonapi.MemberPayload[AdminDetailItem]{Id: record.Id, Type: "pages", Attributes: details}
		default:
			var adminListItem AdminListItem
			json.Unmarshal(payload, &adminListItem)
			result[index] = jsonapi.MemberPayload[AdminListItem]{Id: record.Id, Type: "pages", Attributes: adminListItem}
		}
	}
	return result
}
