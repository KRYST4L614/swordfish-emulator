package v1

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type StorageHandler struct {
	service service.StorageService
}

func NewStorageHandler(service service.StorageService) *StorageHandler {
	return &StorageHandler{
		service: service,
	}
}

func (handler *StorageHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/Storage/{id:[a-zA-Z0-9]+}`, handler.getStorage).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/Storage/{id:[a-zA-Z0-9]+}`, handler.getStorage).Methods(http.MethodGet)
	router.HandleFunc(`/Storage/{id:[a-zA-Z0-9]+}`, handler.updateStorage).Methods(http.MethodPatch)
	router.HandleFunc(`/{root:.*}/Storage/{id:[a-zA-Z0-9]+}`, handler.updateStorage).Methods(http.MethodPatch)
	router.HandleFunc(`/Storage/{id:[a-zA-Z0-9]+}`, handler.replaceStorage).Methods(http.MethodPut)
	router.HandleFunc(`/{root:.*}/Storage/{id:[a-zA-Z0-9]+}`, handler.replaceStorage).Methods(http.MethodPut)
}

func (handler *StorageHandler) getStorage(writer http.ResponseWriter, request *http.Request) {
	serviceRoot, err := handler.service.Get(request.Context(), request.RequestURI)
	if err != nil {
		jsonErr := errlib.GetJSONError(err)
		writer.WriteHeader(jsonErr.Error.Code)
		util.WriteJSON(writer, jsonErr)
		return
	}

	util.WriteJSON(writer, serviceRoot)
}

func (handler *StorageHandler) replaceStorage(writer http.ResponseWriter, request *http.Request) {
	storage, err := util.UnmarshalFromReader[domain.Storage](request.Body)
	if err != nil {
		jsonErr := errlib.GetJSONError(err)
		writer.WriteHeader(jsonErr.Error.Code)
		util.WriteJSON(writer, jsonErr)
		return
	}

	storageId := request.RequestURI
	newStorage, err := handler.service.Replace(request.Context(), storageId, storage)
	if err != nil {
		jsonErr := errlib.GetJSONError(err)
		writer.WriteHeader(jsonErr.Error.Code)
		util.WriteJSON(writer, jsonErr)
		return
	}

	util.WriteJSON(writer, newStorage)
}

func (handler *StorageHandler) updateStorage(writer http.ResponseWriter, request *http.Request) {
	patchData, err := io.ReadAll(request.Body)
	if err != nil {
		jsonErr := errlib.GetJSONError(err)
		writer.WriteHeader(jsonErr.Error.Code)
		util.WriteJSON(writer, jsonErr)
		return
	}

	storageId := request.RequestURI
	newStorage, err := handler.service.Update(request.Context(), storageId, patchData)
	if err != nil {
		jsonErr := errlib.GetJSONError(err)
		writer.WriteHeader(jsonErr.Error.Code)
		util.WriteJSON(writer, jsonErr)
		return
	}

	util.WriteJSON(writer, newStorage)
}
