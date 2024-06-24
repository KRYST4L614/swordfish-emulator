package v1

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
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
	router.HandleFunc(`/{root:.*}/StoragePools`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/StoragePools`+idPathRegex, handler.updateStoragePool).Methods(http.MethodPatch)
	router.HandleFunc(`/{root:.*}/StoragePools`+idPathRegex, handler.replaceStoragePool).Methods(http.MethodPut)
	router.HandleFunc(`/{root:.*}/StoragePools`+idPathRegex, resourceDeleter(handler.service)).Methods(http.MethodDelete)

	router.HandleFunc(`/{root:.*}/ProvidingPools`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/ProvidingPools`+idPathRegex, handler.updateStoragePool).Methods(http.MethodPatch)
	router.HandleFunc(`/{root:.*}/ProvidingPools`+idPathRegex, handler.replaceStoragePool).Methods(http.MethodPut)
	router.HandleFunc(`/{root:.*}/ProvidingPools`+idPathRegex, resourceDeleter(handler.service)).Methods(http.MethodDelete)

	router.HandleFunc(`/{root:.*}/AllocatedPools`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/AllocatedPools`+idPathRegex, handler.updateStoragePool).Methods(http.MethodPatch)
	router.HandleFunc(`/{root:.*}/AllocatedPools`+idPathRegex, handler.replaceStoragePool).Methods(http.MethodPut)
	router.HandleFunc(`/{root:.*}/AllocatedPools`+idPathRegex, resourceDeleter(handler.service)).Methods(http.MethodDelete)
}

func (handler *StoragePoolHandler) replaceStoragePool(writer http.ResponseWriter, request *http.Request) {
	storagePool, err := util.UnmarshalFromReader[domain.StoragePool](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	// TODO: Add validation on incoming resource

	storagePoolId := request.RequestURI
	storagePool.Id = filepath.Base(storagePoolId)
	*storagePool.OdataId = storagePoolId

	newStoragePool, err := handler.service.Replace(request.Context(), storagePoolId, storagePool)
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

	storagePoolId := request.RequestURI
	newStoragePool, err := handler.service.Update(request.Context(), storagePoolId, patchData)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newStoragePool)
}
