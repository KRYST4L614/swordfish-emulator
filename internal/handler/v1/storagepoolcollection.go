package v1

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type StoragePoolCollectionHandler struct {
	service   service.ResourceService
	generator func() (string, error)
}

func NewStoragePoolCollectionHandler(service service.ResourceService) *StoragePoolCollectionHandler {
	return &StoragePoolCollectionHandler{
		service:   service,
		generator: util.IdGenerator(),
	}
}

func (handler *StoragePoolCollectionHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/{root:.*}/ProvidingPools`, handler.getStoragePoolCollection).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/StoragePools`, handler.getStoragePoolCollection).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/AllocatedPools`, handler.getStoragePoolCollection).Methods(http.MethodGet)

	router.HandleFunc(`/{root:.*}/ProvidingPools`, handler.createStoragePool).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/StoragePools`, handler.createStoragePool).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/AllocatedPools`, handler.createStoragePool).Methods(http.MethodPost)

	router.HandleFunc(`/{root:.*}/ProvidingPools/Members`, handler.createStoragePoolFromNotCollectionEndpoint).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/StoragePools/Members`, handler.createStoragePoolFromNotCollectionEndpoint).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/AllocatedPools/Members`, handler.createStoragePoolFromNotCollectionEndpoint).Methods(http.MethodPost)

	router.HandleFunc(`/{root:.*}/ProvidingPools/{id:[a-zA-Z0-9]+}`, handler.createStoragePoolFromNotCollectionEndpoint).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/StoragePools/{id:[a-zA-Z0-9]+}`, handler.createStoragePoolFromNotCollectionEndpoint).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/AllocatedPools/{id:[a-zA-Z0-9]+}`, handler.createStoragePoolFromNotCollectionEndpoint).Methods(http.MethodPost)
}

func (handler *StoragePoolCollectionHandler) getStoragePoolCollection(writer http.ResponseWriter, request *http.Request) {
	serviceRoot, err := handler.service.Get(request.Context(), request.RequestURI)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, serviceRoot)
}

func (handler *StoragePoolCollectionHandler) createStoragePool(writer http.ResponseWriter, request *http.Request) {
	storagePool, err := util.UnmarshalFromReader[domain.StoragePool](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	collectionId := request.RequestURI
	if storagePool.Id == "" {
		storagePool.Id, err = handler.generator()
		if err != nil {
			util.WriteJSONError(writer, err)
			return
		}
	}
	poolId := collectionId + "/" + storagePool.Id
	storagePool.OdataId = util.Addr[string](poolId)

	if _, err := handler.service.Get(request.Context(), collectionId); err != nil {
		_, innerErr := handler.service.CreateCollection(request.Context(), dto.CollectionDto{
			OdataId:   collectionId,
			Name:      "Storage Pool Collection",
			OdataType: "#StoragePoolCollection.StoragePoolCollection",
		})
		if innerErr != nil {
			util.WriteJSONError(writer, err)
			return
		}
	}

	createdStorage, err := handler.service.AddResourceToCollection(request.Context(), collectionId, poolId, storagePool)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	util.WriteJSON(writer, createdStorage)
}

func (handler *StoragePoolCollectionHandler) createStoragePoolFromNotCollectionEndpoint(writer http.ResponseWriter, request *http.Request) {
	if mux.Vars(request)["id"] != "" {
		request = request.WithContext(
			context.WithValue(request.Context(), struct{ Id string }{"Id"}, mux.Vars(request)["id"]),
		)
	}
	request.RequestURI = util.GetParent(request.RequestURI)
	handler.createStoragePool(writer, request)
}
