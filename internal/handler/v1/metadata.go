package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type MetadataHandler struct {
	service service.ResourceService
}

func NewMetadataHandler(service service.ResourceService) *MetadataHandler {
	return &MetadataHandler{service: service}
}

func (handler *MetadataHandler) SetRouter(router *mux.Router) {
	//	router.HandleFunc(`/$metadata`, handler.getMetadataRoot).Methods(http.MethodGet)
	router.HandleFunc(`/odata`, handler.getOData).Methods(http.MethodGet)
}

func (handler *MetadataHandler) getOData(writer http.ResponseWriter, request *http.Request) {
	odata, err := handler.service.Get(request.Context(), request.RequestURI)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, odata)
}
