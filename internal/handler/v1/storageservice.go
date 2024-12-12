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

type StorageServiceHandler struct {
	service service.ResourceService
}

func NewStorageServiceHandler(service service.ResourceService) *StorageServiceHandler {
	return &StorageServiceHandler{
		service: service,
	}
}

func (handler *StorageServiceHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/{root:.*}/StorageServices`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/StorageServices`+idPathRegex, handler.updateStorageService).Methods(http.MethodPatch)
	router.HandleFunc(`/{root:.*}/StorageServices`+idPathRegex, handler.replaceStorageService).Methods(http.MethodPut)
	router.HandleFunc(`/{root:.*}/StorageServices`+idPathRegex, resourceDeleter(handler.service)).Methods(http.MethodDelete)
}

func (handler *StorageServiceHandler) replaceStorageService(writer http.ResponseWriter, request *http.Request) {
	storageService, err := util.UnmarshalFromReader[domain.StorageService](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}


	storageServiceId := request.RequestURI
	storageService.Id = filepath.Base(storageServiceId)
	*storageService.OdataId = storageServiceId

	newStorageService, err := handler.service.Replace(request.Context(), storageServiceId, storageService)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newStorageService)
}

func (handler *StorageServiceHandler) updateStorageService(writer http.ResponseWriter, request *http.Request) {
	patchData, err := io.ReadAll(request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	storageServiceId := request.RequestURI
	newStorageService, err := handler.service.Update(request.Context(), storageServiceId, patchData)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newStorageService)
}
