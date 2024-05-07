package dto

import "gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"

type ResourceDto struct {
	Id   string `db:"id"`
	Data []byte `db:"data"`
}

type CollectionDto struct {
	OdataId   domain.OdataV4Id    `json:"@odata.id"`
	Name      domain.ResourceName `json:"Name"`
	OdataType domain.OdataV4Type  `json:"@odata.type"`
}
