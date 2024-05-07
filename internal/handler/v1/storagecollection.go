package v1

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type StorageCollectionHandler struct {
	service   service.ResourceService
	generator func() (string, error)
}

func NewStorageCollectionHandler(service service.ResourceService) *StorageCollectionHandler {
	return &StorageCollectionHandler{
		service:   service,
		generator: util.IdGenerator(),
	}
}

func (handler *StorageCollectionHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/Storage`, handler.getStorageCollection).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/Storage`, handler.getStorageCollection).Methods(http.MethodGet)

	router.HandleFunc(`/Storage`, handler.createStorage).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/Storage`, handler.createStorage).Methods(http.MethodPost)

	router.HandleFunc(`/Storage/Members`, handler.createStorageFromNotCollectionEndpoint).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/Storage/Members`, handler.createStorageFromNotCollectionEndpoint).Methods(http.MethodPost)

	router.HandleFunc(`/Storage/{id:[a-zA-Z0-9]+}`, handler.createStorageFromNotCollectionEndpoint).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/Storage/{id:[a-zA-Z0-9]+}`, handler.createStorageFromNotCollectionEndpoint).Methods(http.MethodPost)

}

func (handler *StorageCollectionHandler) getStorageCollection(writer http.ResponseWriter, request *http.Request) {
	serviceRoot, err := handler.service.Get(request.Context(), request.RequestURI)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, serviceRoot)
}

func (handler *StorageCollectionHandler) createStorage(writer http.ResponseWriter, request *http.Request) {
	storage, err := util.UnmarshalFromReader[domain.Storage](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	if request.Context().Value("Id") != nil {
		storage.Id = request.Context().Value("Id").(string)
	}

	if storage.Id == "" {
		storage.Id, err = handler.generator()
		if err != nil {
			util.WriteJSONError(writer, err)
			return
		}
	}

	collectionId := request.RequestURI
	storageId := collectionId + "/" + storage.Id
	storage.OdataId = util.Addr[string](storageId)

	createdStorage, err := handler.service.AddResourceToCollection(request.Context(), collectionId, storageId, storage)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	util.WriteJSON(writer, createdStorage)
}

func (handler *StorageCollectionHandler) createStorageFromNotCollectionEndpoint(writer http.ResponseWriter, request *http.Request) {
	if mux.Vars(request)["id"] != "" {
		request = request.WithContext(
			context.WithValue(request.Context(), struct{ Id string }{"Id"}, mux.Vars(request)["id"]),
		)
	}
	request.RequestURI = util.GetParent(request.RequestURI)
	handler.createStorage(writer, request)
}
