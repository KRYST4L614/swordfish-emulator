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

type SystemHandler struct {
	service service.ResourceService
}

func NewSystemHandler(service service.ResourceService) *SystemHandler {
	return &SystemHandler{
		service: service,
	}
}

func (handler *SystemHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/Systems`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/Systems`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)

	router.HandleFunc(`/Systems`+idPathRegex, handler.updateSystem).Methods(http.MethodPatch)
	router.HandleFunc(`/{root:.*}/Systems`+idPathRegex, handler.updateSystem).Methods(http.MethodPatch)

	router.HandleFunc(`/Systems/{id:[a-zA-Z0-9]+}`, handler.replaceSystem).Methods(http.MethodPut)
	router.HandleFunc(`/{root:.*}/Systems/{id:[a-zA-Z0-9]+}`, handler.replaceSystem).Methods(http.MethodPut)

	router.HandleFunc(`/Systems/{id:[a-zA-z0-9]+}`, resourceDeleter(handler.service)).Methods(http.MethodDelete)
	router.HandleFunc(`/{root:.*}/Systems/{id:[a-zA-z0-9]+}`, resourceDeleter(handler.service)).Methods(http.MethodDelete)
}

func (handler *SystemHandler) replaceSystem(writer http.ResponseWriter, request *http.Request) {
	system, err := util.UnmarshalFromReader[domain.ComputerSystem](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	systemId := request.RequestURI
	system.Id = filepath.Base(systemId)
	*system.OdataId = systemId

	newSystem, err := handler.service.Replace(request.Context(), systemId, system)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newSystem)
}

func (handler *SystemHandler) updateSystem(writer http.ResponseWriter, request *http.Request) {
	patchData, err := io.ReadAll(request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	systemId := request.RequestURI
	newSystem, err := handler.service.Update(request.Context(), systemId, patchData)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newSystem)
}
