package jsonapi

type MemberSerializer interface{}

// Schema for jsonapi member response
type memberSerializer[T any] struct {
	Data          interface{}   `json:"data"`
	Included      []interface{} `json:"included"`
	Relationships []interface{} `json:"relationships"`
	Meta          interface{}   `json:"meta"`
}

func NewMemberSerializer[T any](resource interface{}, included []interface{}, relationships []interface{}, meta interface{}) MemberSerializer {
	return &memberSerializer[T]{
		Data:          resource,
		Included:      included,
		Relationships: relationships,
		Meta:          meta,
	}
}

// Collection serializer
type CollectionSerializer[T any] interface{}

// Schema for jsonapi collection response
type collectionSerializer[T any] struct {
	Data []interface{}          `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}

func NewCollectionSerializer[T any](collection []interface{}, meta map[string]interface{}) CollectionSerializer[T] {
	return &collectionSerializer[T]{
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
