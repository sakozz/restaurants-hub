package restaurants

import (
	"encoding/json"
)

func (restaurants Restaurants) Serialize(view ViewTypes) []interface{} {
	result := make([]interface{}, len(restaurants))
	for index, user := range restaurants {
		result[index] = user.Serialize(view)
	}
	return result
}

func (record *Restaurant) Serialize(view ViewTypes) interface{} {
	payload, _ := json.Marshal(record)
	switch view {
	case AdminList:
		var adminListItem AdminListItem
		json.Unmarshal(payload, &adminListItem)
		return Payload[AdminListItem]{Id: record.Id, Type: "restaurants", Attributes: adminListItem}
	default:
		var adminListItem AdminListItem
		json.Unmarshal(payload, &adminListItem)
		return Payload[AdminListItem]{Id: record.Id, Type: "restaurants", Attributes: adminListItem}
	}
}
