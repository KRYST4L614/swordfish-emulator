package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type SystemCollectionHandler struct {
	service service.ResourceService
}

func NewSystemCollectionHandler(service service.ResourceService) *SystemCollectionHandler {
	return &SystemCollectionHandler{
		service: service,
	}
}

func (handler *SystemCollectionHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/Systems`, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/Systems`, resourceGetter(handler.service)).Methods(http.MethodGet)

	router.HandleFunc(`/Systems`, handler.createSystem).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/Systems`, handler.createSystem).Methods(http.MethodPost)

	// Endpoints only for compatibility with official Swordfish Emulator
	//
	// Not a good way, should be removed when Ansible modules will handle
	// resources creation in right way.
	router.HandleFunc(`/Systems`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createSystem)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/Systems`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createSystem)).Methods(http.MethodPost)
}

func (handler *SystemCollectionHandler) createSystem(writer http.ResponseWriter, request *http.Request) {
	system, err := util.UnmarshalFromReader[domain.ComputerSystem](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	// Only for compatibility with official Swordfish Emulator.
	// TODO: Need to be removed when modules handle creation in right way for clearer dataflow.
	contextId := request.Context().Value(idContext)
	if contextId != nil {
		system.Id = contextId.(string)
	}

	createdSystem, err := handler.service.AddResourceToCollection(request.Context(), dto.ResourceRequestDto{
		Name:            system.Name,
		Id:              system.Id,
		OdataType:       "#ComputerSystem.v1_18_0.ComputerSystem",
		Resource:        system,
		IdSetter:        func(id string) { system.Id = id },
		OdataIdSetter:   func(odataId string) { system.OdataId = &odataId },
		OdataTypeSetter: func(odataType string) { system.OdataType = &odataType },
		Collection: dto.CollectionDto{
			OdataId:   request.RequestURI,
			Name:      "Computer System",
			OdataType: "#ComputerSystemCollection.ComputerSystemCollection",
		},
	})

	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	util.WriteJSON(writer, createdSystem)
}
