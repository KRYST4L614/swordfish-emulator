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

type StorageControllerHandler struct {
	service service.ResourceService
}

func NewStorageControllerHandler(service service.ResourceService) *StorageControllerHandler {
	return &StorageControllerHandler{
		service: service,
	}
}

func (handler *StorageControllerHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/Controller`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/Controller`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)

	router.HandleFunc(`/Controller`+idPathRegex, handler.updateStorageController).Methods(http.MethodPatch)
	router.HandleFunc(`/{root:.*}/Controller`+idPathRegex, handler.updateStorageController).Methods(http.MethodPatch)

	router.HandleFunc(`/Controller/{id:[a-zA-Z0-9]+}`, handler.replaceStorageController).Methods(http.MethodPut)
	router.HandleFunc(`/{root:.*}/Controller/{id:[a-zA-Z0-9]+}`, handler.replaceStorageController).Methods(http.MethodPut)

	router.HandleFunc(`/Controller/{id:[a-zA-z0-9]+}`, resourceDeleter(handler.service)).Methods(http.MethodDelete)
	router.HandleFunc(`/{root:.*}/Controller/{id:[a-zA-z0-9]+}`, resourceDeleter(handler.service)).Methods(http.MethodDelete)
}

func (handler *StorageControllerHandler) replaceStorageController(writer http.ResponseWriter, request *http.Request) {
	storage, err := util.UnmarshalFromReader[domain.StorageController](request.Body)
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

func (handler *StorageControllerHandler) updateStorageController(writer http.ResponseWriter, request *http.Request) {
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
