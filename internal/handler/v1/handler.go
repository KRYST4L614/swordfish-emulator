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

// NewHandler - init function to add all implemented handlers
func NewHandler(services *service.Service) *Handler {
	handlers := make([]CommonHandler, 0, 10)
	handlers = append(handlers, NewMetadataHandler(services.ResourceService))
	handlers = append(handlers, NewServiceRootHandler(services.ResourceService))
	handlers = append(handlers, NewSessionHandler(services.ResourceService))
	handlers = append(handlers, NewStorageCollectionHandler(services.ResourceService))
	handlers = append(handlers, NewStorageHandler(services.ResourceService))
	handlers = append(handlers, NewStoragePoolCollectionHandler(services.ResourceService))
	handlers = append(handlers, NewStoragePoolHandler(services.ResourceService))
	handlers = append(handlers, NewStorageServiceCollectionHandler(services.ResourceService))
	handlers = append(handlers, NewStorageServiceHandler(services.ResourceService))
	handlers = append(handlers, NewVolumeCollectionHandler(services.ResourceService))
	handlers = append(handlers, NewVolumeHandler(services.ResourceService))
	handlers = append(handlers, NewFileSystemCollectionHandler(services.ResourceService))
	handlers = append(handlers, NewStorageServiceHandler(services.ResourceService))
	handlers = append(handlers, NewFileSystemCollectionHandler(services.ResourceService))
	handlers = append(handlers, NewFileSystemHandler(services.ResourceService))
	handlers = append(handlers, NewFileShareHandler(services.ResourceService))
	handlers = append(handlers, NewFileShareCollectionHandler(services.ResourceService))
	handlers = append(handlers, NewEthernetInterfaceHandler(services.ResourceService))
	handlers = append(handlers, NewEthernetInterfaceCollectionHandler(services.ResourceService))
	handlers = append(handlers, NewSystemHandler(services.ResourceService))
	handlers = append(handlers, NewSystemCollectionHandler(services.ResourceService))
	handlers = append(handlers, NewFabricCollectionHandler(services.ResourceService))
	handlers = append(handlers, NewFabricHandler(services.ResourceService))
	handlers = append(handlers, NewStorageControllerCollectionHandler(services.ResourceService))
	handlers = append(handlers, NewStorageControllerHandler(services.ResourceService))
	return &Handler{
		handlers: handlers,
	}
}

// SetRouter sets all HandleFuncs for <route>/v1
func (handler *Handler) SetRouter(router *mux.Router) {
	sub := router.PathPrefix("/v1").Subrouter()
	for _, item := range handler.handlers {
		item.SetRouter(sub)
	}
}
