package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	v1 "gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/handler/v1"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/middleware"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type Handler struct {
	v1 *v1.Handler
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		v1: v1.NewHandler(services),
	}
}

func (h *Handler) SetRouter(router *mux.Router) {
	sub := router.PathPrefix("/redfish").Subrouter()
	h.v1.SetRouter(sub)
	router.Use(middleware.Logging)
	router.NotFoundHandler = http.HandlerFunc(h.notFoundHandler)
}

func (h *Handler) notFoundHandler(writer http.ResponseWriter, request *http.Request) {
	logrus.Tracef("Request to not existing path: %s %s", request.Method, request.URL.Path)
	writer.WriteHeader(http.StatusNotFound)
	util.WriteJSONError(writer, errlib.ErrNotFound)
}
