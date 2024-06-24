package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type StorageCollectionHandler struct {
	service service.ResourceService
}

func NewStorageCollectionHandler(service service.ResourceService) *StorageCollectionHandler {
	return &StorageCollectionHandler{
		service: service,
	}
}

func (handler *StorageCollectionHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/Storage`, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/Storage`, resourceGetter(handler.service)).Methods(http.MethodGet)

	router.HandleFunc(`/Storage`, handler.createStorage).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/Storage`, handler.createStorage).Methods(http.MethodPost)

	router.HandleFunc(`/Storage/Members`, resourceCreatorFromNotCollectionEndpoint(handler.createStorage)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/Storage/Members`, resourceCreatorFromNotCollectionEndpoint(handler.createStorage)).Methods(http.MethodPost)

	// Endpoints only for compatibility with official Swordfish Emulator
	//
	// Not a good way, should be removed when Ansible modules will handle
	// resources creation in right way.
	router.HandleFunc(`/Storage`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createStorage)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/Storage`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createStorage)).Methods(http.MethodPost)
}

func (handler *StorageCollectionHandler) createStorage(writer http.ResponseWriter, request *http.Request) {
	storage, err := util.UnmarshalFromReader[domain.Storage](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	// TODO: Add validation on incoming resource

	// Only for compatibility with official Swordfish Emulator.
	// TODO: Need to be removed when modules handle creation in right way for clearer dataflow.
	contextId := request.Context().Value(idContext)
	if contextId != nil {
		storage.Id = contextId.(string)
	}

	createdStorage, err := handler.service.AddResourceToCollection(request.Context(), dto.ResourceRequestDto{
		Name:            storage.Name,
		Id:              storage.Id,
		OdataType:       "#Storage.v1_15_1.Storage",
		Resource:        storage,
		IdSetter:        func(id string) { storage.Id = id },
		OdataIdSetter:   func(odataId string) { storage.OdataId = &odataId },
		OdataTypeSetter: func(odataType string) { storage.OdataType = &odataType },
		Collection: dto.CollectionDto{
			OdataId:   request.RequestURI,
			Name:      "Storage Collection",
			OdataType: "#StorageCollection.StorageCollection",
		},
	})

	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	util.WriteJSON(writer, createdStorage)
}
