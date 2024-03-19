package v1

// import (
// 	"github.com/gorilla/mux"
// 	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
// )

// type Handler struct {
// 	handlers []ResourceHandler
// }

// type ResourceHandler interface {
// 	SetRouter(router *mux.Router)
// }

// func NewHandler(services *service.Service) *Handler {
// 	handlers := make([]ResourceHandler, 0)
// 	handlers = append(handlers, NewVolumeHandler(services.VolumeService))
// 	return &Handler{
// 		handlers: handlers,
// 	}
// }

// func (handler *Handler) SetRouter(router *mux.Router) {
// 	sub := router.PathPrefix("/v1").Subrouter()
// 	for _, item := range handler.handlers {
// 		item.SetRouter(sub)
// 	}
// }
