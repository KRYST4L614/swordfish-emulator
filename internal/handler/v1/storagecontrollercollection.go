package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type StorageControllerCollectionHandler struct {
	service service.ResourceService
}

func NewStorageControllerCollectionHandler(service service.ResourceService) *StorageControllerCollectionHandler {
	return &StorageControllerCollectionHandler{
		service: service,
	}
}

func (handler *StorageControllerCollectionHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/Controller`, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/Controller`, resourceGetter(handler.service)).Methods(http.MethodGet)

	router.HandleFunc(`/Controller`, handler.createStorageController).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/Controller`, handler.createStorageController).Methods(http.MethodPost)

	router.HandleFunc(`/Controller/Members`, resourceCreatorFromNotCollectionEndpoint(handler.createStorageController)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/Controller/Members`, resourceCreatorFromNotCollectionEndpoint(handler.createStorageController)).Methods(http.MethodPost)

	// Endpoints only for compatibility with official Swordfish Emulator
	//
	// Not a good way, should be removed when Ansible modules will handle
	// resources creation in right way.
	router.HandleFunc(`/Controller`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createStorageController)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/Controller`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createStorageController)).Methods(http.MethodPost)
}

func (handler *StorageControllerCollectionHandler) createStorageController(writer http.ResponseWriter, request *http.Request) {
	storage, err := util.UnmarshalFromReader[domain.StorageController](request.Body)
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
		OdataType:       "#StorageController.v1_6_0.StorageController",
		Resource:        storage,
		IdSetter:        func(id string) { storage.Id = id },
		OdataIdSetter:   func(odataId string) { storage.OdataId = &odataId },
		OdataTypeSetter: func(odataType string) { storage.OdataType = &odataType },
		Collection: dto.CollectionDto{
			OdataId:   request.RequestURI,
			Name:      "Storage Controller Collection",
			OdataType: "#StorageControllerCollection.StorageControllerCollection",
		},
	})

	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	util.WriteJSON(writer, createdStorage)
}
