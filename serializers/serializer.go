package serializers

// Schema for jsonapi member response
type MemberSerializer struct {
	Data          interface{}   `json:"data"`
	Included      []interface{} `json:"included"`
	Relationships []interface{} `json:"relationships"`
	Meta          interface{}   `json:"meta"`
}

func NewMemberSerializer(resource interface{}, included []interface{}, relationships []interface{}, meta interface{}) *MemberSerializer {
	return &MemberSerializer{
		Data:          resource,
		Included:      included,
		Relationships: relationships,
		Meta:          meta,
	}
}

// Schema for jsonapi collection response
type CollectionSerializer struct {
	Data []interface{}          `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}

func NewCollectionSerializer(collection []interface{}, meta map[string]interface{}) *CollectionSerializer {
	return &CollectionSerializer{
		Data: collection,
		Meta: meta,
	}
}

// Schema for jsonapi resource
type MemberPayload[T any] struct {
	Id         int64  `json:"id"`
	Type       string `json:"type"`
	Attributes T      `json:"attributes"`
}
