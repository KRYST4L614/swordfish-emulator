package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type StorageCollectionHandler struct {
	service service.StorageCollectionService
}

func NewStorageCollectionHandler(service service.StorageCollectionService) *StorageCollectionHandler {
	return &StorageCollectionHandler{
		service: service,
	}
}

func (handler *StorageCollectionHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/Storage`, handler.getStorageCollection).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/Storage`, handler.getStorageCollection).Methods(http.MethodGet)
	router.HandleFunc(`/Storage`, handler.createStorage).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/Storage`, handler.createStorage).Methods(http.MethodPost)
	router.HandleFunc(`/Storage/{id:[a-zA-z0-9]+}`, handler.deleteStorage).Methods(http.MethodDelete)
	router.HandleFunc(`/{root:.*}/Storage/{id:[a-zA-z0-9]+}`, handler.deleteStorage).Methods(http.MethodDelete)
}

func (handler *StorageCollectionHandler) getStorageCollection(writer http.ResponseWriter, request *http.Request) {
	serviceRoot, err := handler.service.Get(request.Context(), request.RequestURI)
	if err != nil {
		jsonErr := errlib.GetJSONError(err)
		writer.WriteHeader(jsonErr.Error.Code)
		util.WriteJSON(writer, jsonErr)
		return
	}

	util.WriteJSON(writer, serviceRoot)
}

func (handler *StorageCollectionHandler) createStorage(writer http.ResponseWriter, request *http.Request) {
	storage, err := util.UnmarshalFromReader[domain.Storage](request.Body)
	if err != nil {
		jsonErr := errlib.GetJSONError(err)
		writer.WriteHeader(jsonErr.Error.Code)
		util.WriteJSON(writer, jsonErr)
		return
	}

	err = handler.service.AddStorage(request.Context(), request.RequestURI, storage)
	if err != nil {
		jsonErr := errlib.GetJSONError(err)
		writer.WriteHeader(jsonErr.Error.Code)
		util.WriteJSON(writer, jsonErr)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	util.WriteJSON(writer, storage)
}

func (handler *StorageCollectionHandler) deleteStorage(writer http.ResponseWriter, request *http.Request) {
	storageId := request.RequestURI
	storage, err := handler.service.DeleteStorage(request.Context(), storageId)
	if err != nil {
		jsonErr := errlib.GetJSONError(err)
		writer.WriteHeader(jsonErr.Error.Code)
		util.WriteJSON(writer, jsonErr)
		return
	}
	writer.WriteHeader(http.StatusOK)
	util.WriteJSON(writer, storage)
}
