package v1

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type EthernetInterfaceHandler struct {
	service service.ResourceService
}

func NewEthernetInterfaceHandler(service service.ResourceService) *EthernetInterfaceHandler {
	return &EthernetInterfaceHandler{
		service: service,
	}
}

func (handler *EthernetInterfaceHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/EthernetInterfaces`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/EthernetInterfaces`+idPathRegex, handler.updateEthernetInterface).Methods(http.MethodPatch)
	router.HandleFunc(`/EthernetInterfaces`+idPathRegex, handler.replaceEthernetInterface).Methods(http.MethodPut)
	router.HandleFunc(`/EthernetInterfaces`+idPathRegex, resourceDeleter(handler.service)).Methods(http.MethodDelete)

	router.HandleFunc(`/{root:.*}/EthernetInterfaces`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/EthernetInterfaces`+idPathRegex, handler.updateEthernetInterface).Methods(http.MethodPatch)
	router.HandleFunc(`/{root:.*}/EthernetInterfaces`+idPathRegex, handler.replaceEthernetInterface).Methods(http.MethodPut)
	router.HandleFunc(`/{root:.*}/EthernetInterfaces`+idPathRegex, resourceDeleter(handler.service)).Methods(http.MethodDelete)
}

func (handler *EthernetInterfaceHandler) replaceEthernetInterface(writer http.ResponseWriter, request *http.Request) {
	ethernetInterface, err := util.UnmarshalFromReader[domain.EthernetInterface](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	ethernetInterfaceId := request.RequestURI
	ethernetInterface.Id = filepath.Base(ethernetInterfaceId)
	*ethernetInterface.OdataId = ethernetInterfaceId

	newEthernetInterface, err := handler.service.Replace(request.Context(), ethernetInterfaceId, ethernetInterface)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newEthernetInterface)
}

func (handler *EthernetInterfaceHandler) updateEthernetInterface(writer http.ResponseWriter, request *http.Request) {
	patchData, err := io.ReadAll(request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	ethernetInterfaceId := request.RequestURI
	newEthernetInterface, err := handler.service.Update(request.Context(), ethernetInterfaceId, patchData)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newEthernetInterface)
}
