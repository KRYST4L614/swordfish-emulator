package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
)

type MetadataHandler struct {
	service service.ResourceService
}

func NewMetadataHandler(service service.ResourceService) *MetadataHandler {
	return &MetadataHandler{service: service}
}

func (handler *MetadataHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/$metadata`, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/odata`, resourceGetter(handler.service)).Methods(http.MethodGet)
}
