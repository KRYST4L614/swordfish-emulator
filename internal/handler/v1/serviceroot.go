package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
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
	router.HandleFunc("", resourceGetter(handler.service)).Methods(http.MethodGet)
}
