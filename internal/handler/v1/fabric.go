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

type FabricHandler struct {
	service service.ResourceService
}

func NewFabricHandler(service service.ResourceService) *FabricHandler {
	return &FabricHandler{
		service: service,
	}
}

func (handler *FabricHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/Fabric`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/Fabric`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)

	router.HandleFunc(`/Fabric`+idPathRegex, handler.updateFabric).Methods(http.MethodPatch)
	router.HandleFunc(`/{root:.*}/Fabric`+idPathRegex, handler.updateFabric).Methods(http.MethodPatch)

	router.HandleFunc(`/Fabric/{id:[a-zA-Z0-9]+}`, handler.replaceFabric).Methods(http.MethodPut)
	router.HandleFunc(`/{root:.*}/Fabric/{id:[a-zA-Z0-9]+}`, handler.replaceFabric).Methods(http.MethodPut)

	router.HandleFunc(`/Fabric/{id:[a-zA-z0-9]+}`, resourceDeleter(handler.service)).Methods(http.MethodDelete)
	router.HandleFunc(`/{root:.*}/Fabric/{id:[a-zA-z0-9]+}`, resourceDeleter(handler.service)).Methods(http.MethodDelete)
}

func (handler *FabricHandler) replaceFabric(writer http.ResponseWriter, request *http.Request) {
	storage, err := util.UnmarshalFromReader[domain.Fabric](request.Body)
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

func (handler *FabricHandler) updateFabric(writer http.ResponseWriter, request *http.Request) {
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
