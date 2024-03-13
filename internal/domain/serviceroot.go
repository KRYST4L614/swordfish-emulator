package domain

type ServiceRoot struct {
	Base           `json:",inline"`
	RedfishVersion string        `json:"RedfishVersion"`
	UUID           string        `json:"UUID"`
	Chassis        InlineODataId `json:"Chassis"`
	Fabrics        InlineODataId `json:"Fabrics"`
	Managers       InlineODataId `json:"Managers"`
	SessionService InlineODataId `json:"SessionService"`
	Registries     InlineODataId `json:"Registries"`
	Storage        InlineODataId `json:"Storage"`
}
