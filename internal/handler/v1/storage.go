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

type StorageHandler struct {
	service service.ResourceService
}

func NewStorageHandler(service service.ResourceService) *StorageHandler {
	return &StorageHandler{
		service: service,
	}
}

func (handler *StorageHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/Storage`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/Storage`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)

	router.HandleFunc(`/Storage`+idPathRegex, handler.updateStorage).Methods(http.MethodPatch)
	router.HandleFunc(`/{root:.*}/Storage`+idPathRegex, handler.updateStorage).Methods(http.MethodPatch)

	router.HandleFunc(`/Storage/{id:[a-zA-Z0-9]+}`, handler.replaceStorage).Methods(http.MethodPut)
	router.HandleFunc(`/{root:.*}/Storage/{id:[a-zA-Z0-9]+}`, handler.replaceStorage).Methods(http.MethodPut)

	router.HandleFunc(`/Storage/{id:[a-zA-z0-9]+}`, resourceDeleter(handler.service)).Methods(http.MethodDelete)
	router.HandleFunc(`/{root:.*}/Storage/{id:[a-zA-z0-9]+}`, resourceDeleter(handler.service)).Methods(http.MethodDelete)
}

func (handler *StorageHandler) replaceStorage(writer http.ResponseWriter, request *http.Request) {
	storage, err := util.UnmarshalFromReader[domain.Storage](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	storageId := request.RequestURI
	storage.Id = filepath.Base(storageId)
	*storage.OdataId = storageId

	newStorage, err := handler.service.Replace(request.Context(), storageId, storage)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newStorage)
}

func (handler *StorageHandler) updateStorage(writer http.ResponseWriter, request *http.Request) {
	patchData, err := io.ReadAll(request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	storageId := request.RequestURI
	newStorage, err := handler.service.Update(request.Context(), storageId, patchData)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newStorage)
}
