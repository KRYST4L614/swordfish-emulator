package v1

import (
	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
)

type Handler struct {
	handlers []CommonHandler
}

type CommonHandler interface {
	SetRouter(router *mux.Router)
}

func NewHandler(services *service.Service) *Handler {
	handlers := make([]CommonHandler, 0)
	handlers = append(handlers, NewServiceRootHandler(services.ResourceService))
	handlers = append(handlers, NewStorageCollectionHandler(services.ResourceService))
	handlers = append(handlers, NewStorageHandler(services.ResourceService))
	handlers = append(handlers, NewStoragePoolCollectionHandler(services.ResourceService))
	handlers = append(handlers, NewStoragePoolHandler(services.ResourceService))
	return &Handler{
		handlers: handlers,
	}
}

func (handler *Handler) SetRouter(router *mux.Router) {
	sub := router.PathPrefix("/v1").Subrouter()
	for _, item := range handler.handlers {
		item.SetRouter(sub)
	}
}
