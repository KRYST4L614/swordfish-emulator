package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type StoragePoolCollectionHandler struct {
	service service.ResourceService
}

func NewStoragePoolCollectionHandler(service service.ResourceService) *StoragePoolCollectionHandler {
	return &StoragePoolCollectionHandler{
		service: service,
	}
}

func (handler *StoragePoolCollectionHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/{root:.*}/ProvidingPools`, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/StoragePools`, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/AllocatedPools`, resourceGetter(handler.service)).Methods(http.MethodGet)

	router.HandleFunc(`/{root:.*}/ProvidingPools`, handler.createStoragePool).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/StoragePools`, handler.createStoragePool).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/AllocatedPools`, handler.createStoragePool).Methods(http.MethodPost)

	router.HandleFunc(`/{root:.*}/ProvidingPools/Members`, resourceCreatorFromNotCollectionEndpoint(handler.createStoragePool)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/StoragePools/Members`, resourceCreatorFromNotCollectionEndpoint(handler.createStoragePool)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/AllocatedPools/Members`, resourceCreatorFromNotCollectionEndpoint(handler.createStoragePool)).Methods(http.MethodPost)

	// Endpoints only for compatibility with official Swordfish Emulator
	//
	// Not a good way, should be removed when Ansible modules will handle
	// resources creation in right way
	router.HandleFunc(`/{root:.*}/ProvidingPools`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createStoragePool)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/StoragePools`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createStoragePool)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/AllocatedPools`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createStoragePool)).Methods(http.MethodPost)
}

func (handler *StoragePoolCollectionHandler) createStoragePool(writer http.ResponseWriter, request *http.Request) {
	storagePool, err := util.UnmarshalFromReader[domain.StoragePool](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	// Only for compatibility with official Swordfish Emulator.
	// TODO: Need to be removed when modules handle creation in right way for clearer dataflow.
	contextId := request.Context().Value(idContext)
	if contextId != nil {
		storagePool.Id = contextId.(string)
	}

	// TODO: Add validation on incoming resource
	logrus.Info("Pool uri: ", request.RequestURI)
	pool, err := handler.service.AddResourceToCollection(request.Context(), dto.ResourceRequestDto{
		Name:            storagePool.Name,
		Id:              storagePool.Id,
		OdataType:       "#StoragePool.v1_9_0.StoragePool",
		Resource:        storagePool,
		IdSetter:        func(id string) { storagePool.Id = id },
		OdataIdSetter:   func(odataId string) { storagePool.OdataId = &odataId },
		OdataTypeSetter: func(odataType string) { storagePool.OdataType = &odataType },
		Collection: dto.CollectionDto{
			OdataId:   request.RequestURI,
			Name:      "Storage Pool Collection",
			OdataType: "#StoragePoolCollection.StoragePoolCollection",
		},
	})

	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	util.WriteJSON(writer, pool)
}
