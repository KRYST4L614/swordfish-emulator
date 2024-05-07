package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

// ServiceRootHandler addresses to the ServiceRoot endpoint
type ServiceRootHandler struct {
	service service.ResourceService
}

func NewServiceRootHandler(service service.ResourceService) *ServiceRootHandler {
	return &ServiceRootHandler{
		service: service,
	}
}

// SetRouter sets handle functions for operation on ServiceRoot resource
func (handler *ServiceRootHandler) SetRouter(router *mux.Router) {
	router.HandleFunc("", handler.getServiceRoot).Methods(http.MethodGet)
}

func (handler *ServiceRootHandler) getServiceRoot(writer http.ResponseWriter, request *http.Request) {
	serviceRoot, err := handler.service.Get(request.Context(), request.RequestURI)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, serviceRoot)
}
