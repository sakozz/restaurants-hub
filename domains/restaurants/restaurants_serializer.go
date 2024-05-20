package restaurants

import (
	"encoding/json"

	"resturants-hub.com/m/v2/jsonapi"
)

func (record *Restaurant) MemberFor(view ViewTypes) interface{} {
	payload, _ := json.Marshal(record)
	switch view {
	case AdminList:
		var adminListItem AdminListItem
		json.Unmarshal(payload, &adminListItem)
		return jsonapi.MemberPayload[AdminListItem]{Id: record.Id, Type: "restaurants", Attributes: adminListItem}
	case AdminDetails:
		var details AdminDetailItem
		json.Unmarshal(payload, &details)
		return jsonapi.MemberPayload[AdminDetailItem]{Id: record.Id, Type: "restaurants", Attributes: details}
	default:
		var adminListItem AdminListItem
		json.Unmarshal(payload, &adminListItem)
		return jsonapi.MemberPayload[AdminListItem]{Id: record.Id, Type: "restaurants", Attributes: adminListItem}
	}
}

func (restaurants Restaurants) CollectionFor(view ViewTypes) []interface{} {
	result := make([]interface{}, len(restaurants))
	for index, record := range restaurants {
		payload, _ := json.Marshal(record)
		switch view {
		case AdminList:
			var adminListItem AdminListItem
			json.Unmarshal(payload, &adminListItem)
			result[index] = jsonapi.MemberPayload[AdminListItem]{Id: record.Id, Type: "restaurants", Attributes: adminListItem}
		case AdminDetails:
			var details AdminDetailItem
			json.Unmarshal(payload, &details)
			result[index] = jsonapi.MemberPayload[AdminDetailItem]{Id: record.Id, Type: "restaurants", Attributes: details}
		default:
			var adminListItem AdminListItem
			json.Unmarshal(payload, &adminListItem)
			result[index] = jsonapi.MemberPayload[AdminListItem]{Id: record.Id, Type: "restaurants", Attributes: adminListItem}
		}
	}
	return result
}
