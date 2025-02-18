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

type StorageServiceCollectionHandler struct {
	service service.ResourceService
}

func NewStorageServiceCollectionHandler(service service.ResourceService) *StorageServiceCollectionHandler {
	return &StorageServiceCollectionHandler{
		service: service,
	}
}

func (handler *StorageServiceCollectionHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/StorageServices`, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/StorageServices`, resourceGetter(handler.service)).Methods(http.MethodGet)

	router.HandleFunc(`/StorageServices`, handler.createStorageService).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/StorageServices`, handler.createStorageService).Methods(http.MethodPost)

	router.HandleFunc(`/StorageServices`, resourceCreatorFromNotCollectionEndpoint(handler.createStorageService)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/StorageServices`, resourceCreatorFromNotCollectionEndpoint(handler.createStorageService)).Methods(http.MethodPost)

	// Endpoints only for compatibility with official Swordfish Emulator
	//
	// Not a good way, should be removed when Ansible modules will handle
	// resources creation in right way
	router.HandleFunc(`/StorageServices`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createStorageService)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/StorageServices`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createStorageService)).Methods(http.MethodPost)
}

func (handler *StorageServiceCollectionHandler) createStorageService(writer http.ResponseWriter, request *http.Request) {
	storageService, err := util.UnmarshalFromReader[domain.StorageService](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	// Only for compatibility with official Swordfish Emulator.
	// TODO: Need to be removed when modules handle creation in right way for clearer dataflow.

	slog.Info("StorageServieCollection uri: " + request.RequestURI)
	createdStorageService, err := handler.service.AddResourceToCollection(request.Context(), dto.ResourceRequestDto{
		Name:            storageService.Name,
		Id:              storageService.Id,
		OdataType:       "#StorageServiceCollection.StorageServiceCollection",
		Resource:        storageService,
		IdSetter:        func(id string) { storageService.Id = id },
		OdataIdSetter:   func(odataId string) { storageService.OdataId = &odataId },
		OdataTypeSetter: func(odataType string) { storageService.OdataType = &odataType },
		Collection: dto.CollectionDto{
			OdataId:   request.RequestURI,
			Name:      "StorageServiceCollection",
			OdataType: "#StorageService.v1_7_0.StorageService",
		},
	})

	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	util.WriteJSON(writer, createdStorageService)
}
