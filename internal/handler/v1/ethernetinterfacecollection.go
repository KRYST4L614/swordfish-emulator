package v1

import (
	"net/http"

	"log/slog"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type EthernetInterfaceCollectionHandler struct {
	service service.ResourceService
}

func NewEthernetInterfaceCollectionHandler(service service.ResourceService) *EthernetInterfaceCollectionHandler {
	return &EthernetInterfaceCollectionHandler{
		service: service,
	}
}

func (handler *EthernetInterfaceCollectionHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/EthernetInterfaces`, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/EthernetInterfaces`, resourceGetter(handler.service)).Methods(http.MethodGet)

	router.HandleFunc(`/EthernetInterfaces`, handler.createEthernetInterface).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/EthernetInterfaces`, handler.createEthernetInterface).Methods(http.MethodPost)

	router.HandleFunc(`/EthernetInterfaces`, resourceCreatorFromNotCollectionEndpoint(handler.createEthernetInterface)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/EthernetInterfaces`, resourceCreatorFromNotCollectionEndpoint(handler.createEthernetInterface)).Methods(http.MethodPost)

	// Endpoints only for compatibility with official Swordfish Emulator
	//
	// Not a good way, should be removed when Ansible modules will handle
	// resources creation in right way
	router.HandleFunc(`/EthernetInterfaces`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createEthernetInterface)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/EthernetInterfaces`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createEthernetInterface)).Methods(http.MethodPost)
}

func (handler *EthernetInterfaceCollectionHandler) createEthernetInterface(writer http.ResponseWriter, request *http.Request) {
	ethernetInterface, err := util.UnmarshalFromReader[domain.EthernetInterface](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	// Only for compatibility with official Swordfish Emulator.
	// TODO: Need to be removed when modules handle creation in right way for clearer dataflow.

	slog.Info("EthernetInterfaceCollection uri: " + request.RequestURI)
	createdEthernetInterface, err := handler.service.AddResourceToCollection(request.Context(), dto.ResourceRequestDto{
		Name:            ethernetInterface.Name,
		Id:              ethernetInterface.Id,
		OdataType:       "#EthernetInterface.v1_7_0.EthernetInterface",
		Resource:        ethernetInterface,
		IdSetter:        func(id string) { ethernetInterface.Id = id },
		OdataIdSetter:   func(odataId string) { ethernetInterface.OdataId = &odataId },
		OdataTypeSetter: func(odataType string) { ethernetInterface.OdataType = &odataType },
		Collection: dto.CollectionDto{
			OdataId:   request.RequestURI,
			Name:      "EthernetInterfaceCollection",
			OdataType: "#EthernetInterfaceCollection.EthernetInterfaceCollection",
		},
	})

	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	util.WriteJSON(writer, createdEthernetInterface)
}
