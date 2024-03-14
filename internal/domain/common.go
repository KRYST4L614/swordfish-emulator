package domain

type InlineODataId struct {
	ODataId string `json:"@odata.id"`
}

type InlineLinks struct {
	Links map[string]interface{} `json:"Links"`
}

type Base struct {
	ODataContext  string `json:"@odata.context,omitempty"`
	ODataEtag     string `json:"@odata.etag,omitempty"`
	InlineODataId `json:",inline"`
	ODataType     string `json:"@odata.type"`
	Description   string `json:"Description,omitempty"`
	Id            string `json:"Id"`
	Name          string `json:"Name"`
	InlineLinks   `json:",inline,omitempty"`
}
