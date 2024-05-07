package v1

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type StoragePoolHandler struct {
	service service.ResourceService
}

func NewStoragePoolHandler(service service.ResourceService) *StoragePoolHandler {
	return &StoragePoolHandler{
		service: service,
	}
}

func (handler *StoragePoolHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/{root:.*}/StoragePools/{id:[a-zA-Z0-9]+}`, handler.getStoragePool).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/StoragePools/{id:[a-zA-Z0-9]+}`, handler.updateStoragePool).Methods(http.MethodPatch)
	router.HandleFunc(`/{root:.*}/StoragePools/{id:[a-zA-Z0-9]+}`, handler.replaceStoragePool).Methods(http.MethodPut)
	router.HandleFunc(`/{root:.*}/StoragePools/{id:[a-zA-z0-9]+}`, handler.deleteStoragePool).Methods(http.MethodDelete)

	router.HandleFunc(`/{root:.*}/ProvidingPools/{id:[a-zA-Z0-9]+}`, handler.getStoragePool).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/ProvidingPools/{id:[a-zA-Z0-9]+}`, handler.updateStoragePool).Methods(http.MethodPatch)
	router.HandleFunc(`/{root:.*}/ProvidingPools/{id:[a-zA-Z0-9]+}`, handler.replaceStoragePool).Methods(http.MethodPut)
	router.HandleFunc(`/{root:.*}/ProvidingPools/{id:[a-zA-z0-9]+}`, handler.deleteStoragePool).Methods(http.MethodDelete)

	router.HandleFunc(`/{root:.*}/AllocatedPools/{id:[a-zA-Z0-9]+}`, handler.getStoragePool).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/AllocatedPools/{id:[a-zA-Z0-9]+}`, handler.updateStoragePool).Methods(http.MethodPatch)
	router.HandleFunc(`/{root:.*}/AllocatedPools/{id:[a-zA-Z0-9]+}`, handler.replaceStoragePool).Methods(http.MethodPut)
	router.HandleFunc(`/{root:.*}/AllocatedPools/{id:[a-zA-z0-9]+}`, handler.deleteStoragePool).Methods(http.MethodDelete)
}

func (handler *StoragePoolHandler) getStoragePool(writer http.ResponseWriter, request *http.Request) {
	serviceRoot, err := handler.service.Get(request.Context(), request.RequestURI)
	if err != nil {
		jsonErr := errlib.GetJSONError(err)
		writer.WriteHeader(jsonErr.Error.Code)
		util.WriteJSON(writer, jsonErr)
		return
	}

	util.WriteJSON(writer, serviceRoot)
}

func (handler *StoragePoolHandler) replaceStoragePool(writer http.ResponseWriter, request *http.Request) {
	StoragePool, err := util.UnmarshalFromReader[domain.StoragePool](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	StoragePoolId := request.RequestURI
	StoragePool.Id = filepath.Base(StoragePoolId)
	*StoragePool.OdataId = StoragePoolId

	newStoragePool, err := handler.service.Replace(request.Context(), StoragePoolId, StoragePool)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newStoragePool)
}

func (handler *StoragePoolHandler) updateStoragePool(writer http.ResponseWriter, request *http.Request) {
	patchData, err := io.ReadAll(request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	StoragePoolId := request.RequestURI
	newStoragePool, err := handler.service.Update(request.Context(), StoragePoolId, patchData)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newStoragePool)
}

func (handler *StoragePoolHandler) deleteStoragePool(writer http.ResponseWriter, request *http.Request) {
	StoragePoolId := request.RequestURI
	StoragePool, err := handler.service.DeleteResourceFromCollection(request.Context(), util.GetParent(StoragePoolId), StoragePoolId)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusOK)
	util.WriteJSON(writer, StoragePool)
}
