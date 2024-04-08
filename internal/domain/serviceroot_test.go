package domain

// import (
// 	"encoding/json"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestServiceRootValidDeserialization(t *testing.T) {
// 	serviceRootString := `
// 	{
// 		"@odata.type": "#ServiceRoot.v1_15_0.ServiceRoot",
// 		"Id": "RootService",
// 		"Name": "Root Service",
// 		"RedfishVersion": "1.18.0",
// 		"UUID": "92384634-2938-2342-8820-489239905423",
// 		"Chassis": {
// 			"@odata.id": "/redfish/v1/Chassis"
// 		},
// 		"Fabrics": {
// 			"@odata.id": "/redfish/v1/Fabrics"
// 		},
// 		"Managers": {
// 			"@odata.id": "/redfish/v1/Managers"
// 		},
// 		"SessionService": {
// 			"@odata.id": "/redfish/v1/SessionService"
// 		},
// 		"Registries": {
// 			"@odata.id": "/redfish/v1/Registries"
// 		},
// 		"Storage": {
// 			"@odata.id": "/redfish/v1/Storage"
// 		},
// 		"Links": {
// 			"Sessions": {
// 				"@odata.id": "/redfish/v1/SessionService/Sessions"
// 			}
// 		},
// 		"@odata.id": "/redfish/v1"
// 	}
// 	`

// 	var serviceRoot ServiceRoot
// 	err := json.Unmarshal([]byte(serviceRootString), &serviceRoot)
// 	if !assert.NoError(t, err) {
// 		assert.FailNow(t, err.Error())
// 	}

// 	assert.Equal(t, serviceRoot.Id, "RootService")

// 	links := serviceRoot.Links
// 	sessionsRef := links["Sessions"]

// 	sessionsBytes, err := json.Marshal(sessionsRef)
// 	if !assert.NoError(t, err) {
// 		assert.FailNow(t, err.Error())
// 	}

// 	var odata InlineODataId
// 	err = json.Unmarshal(sessionsBytes, &odata)
// 	if !assert.NoError(t, err) {
// 		assert.FailNow(t, err.Error())
// 	}

// 	assert.Equal(t, odata.ODataId, "/redfish/v1/SessionService/Sessions")
// }
