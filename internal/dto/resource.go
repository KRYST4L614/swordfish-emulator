package dto

import "gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"

type ResourceDto struct {
	Id   string `db:"id"`
	Data []byte `db:"data"`
}

// CollectionDto - struct with collection fields, that are
// common for all Swordfish collections
//
// Can be used when you don't now what collection is stored in bytes,
// but want to know some important details after deserialization
type CollectionDto struct {
	OdataId   domain.OdataV4Id    `json:"@odata.id"`
	Name      domain.ResourceName `json:"Name"`
	OdataType domain.OdataV4Type  `json:"@odata.type"`
}

type ResourceRequestDto struct {
	Name      domain.ResourceName
	Id        domain.ResourceId
	OdataType domain.OdataV4Type

	// Resource - field containing pointer to original resource struct,
	// deserialized from request
	Resource interface{}

	// IdSetter - closure, containing original resource struct, that sets
	// Id field for a resource. Basically sets Id for a 'Resource' in this struct.
	IdSetter func(id string)

	// OdataIdSetter - closure, containing original resource struct.
	// Basically sets OdataId for a 'Resource' in this struct.
	OdataIdSetter func(odataId string)

	// OdataTypeSetter - closure, containing original resource struct.
	// Basically sets OdataType for a 'Resource' in this struct.
	OdataTypeSetter func(odataType string)

	// Field should contain details about collection, containing underlying
	// resource.
	Collection CollectionDto
}
