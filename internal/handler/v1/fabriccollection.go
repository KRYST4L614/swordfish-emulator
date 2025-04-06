package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type FabricCollectionHandler struct {
	service service.ResourceService
}

func NewFabricCollectionHandler(service service.ResourceService) *FabricCollectionHandler {
	return &FabricCollectionHandler{
		service: service,
	}
}

func (handler *FabricCollectionHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/Fabric`, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/Fabric`, resourceGetter(handler.service)).Methods(http.MethodGet)

	router.HandleFunc(`/Fabric`, handler.createFabric).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/Fabric`, handler.createFabric).Methods(http.MethodPost)

	router.HandleFunc(`/Fabric/Members`, resourceCreatorFromNotCollectionEndpoint(handler.createFabric)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/Fabric/Members`, resourceCreatorFromNotCollectionEndpoint(handler.createFabric)).Methods(http.MethodPost)

	// Endpoints only for compatibility with official Swordfish Emulator
	//
	// Not a good way, should be removed when Ansible modules will handle
	// resources creation in right way.
	router.HandleFunc(`/Fabric`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createFabric)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/Fabric`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createFabric)).Methods(http.MethodPost)
}

func (handler *FabricCollectionHandler) createFabric(writer http.ResponseWriter, request *http.Request) {
	storage, err := util.UnmarshalFromReader[domain.Storage](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	// Only for compatibility with official Swordfish Emulator.
	// TODO: Need to be removed when modules handle creation in right way for clearer dataflow.
	contextId := request.Context().Value(idContext)
	if contextId != nil {
		storage.Id = contextId.(string)
	}

	createdStorage, err := handler.service.AddResourceToCollection(request.Context(), dto.ResourceRequestDto{
		Name:            storage.Name,
		Id:              storage.Id,
		OdataType:       "#Fabric.v1_3_2.Fabric",
		Resource:        storage,
		IdSetter:        func(id string) { storage.Id = id },
		OdataIdSetter:   func(odataId string) { storage.OdataId = &odataId },
		OdataTypeSetter: func(odataType string) { storage.OdataType = &odataType },
		Collection: dto.CollectionDto{
			OdataId:   request.RequestURI,
			Name:      "Fabric Collection",
			OdataType: "#FabricCollection.FabricCollection",
		},
	})

	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	util.WriteJSON(writer, createdStorage)
}
