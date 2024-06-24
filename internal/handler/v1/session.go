package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

// TODO: this handler is not complete, just stub
type SessionHandler struct {
	service   service.ResourceService
	generator func() (string, error)
}

func NewSessionHandler(service service.ResourceService) *SessionHandler {
	return &SessionHandler{
		service:   service,
		generator: util.IdGenerator(),
	}
}

func (handler *SessionHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/SessionService/Sessions{any:.*}`, handler.makeSession).Methods(http.MethodPost, http.MethodGet, http.MethodDelete)
}

func (handler *SessionHandler) makeSession(writer http.ResponseWriter, request *http.Request) {
	sessionToken, err := handler.generator()
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	// TODO: this is just stub for testing on Ansible modules
	writer.Header().Add("X-Auth-Token", sessionToken)
	writer.Header().Add("Location", request.RequestURI+"/"+sessionToken)
	writer.Header().Add("Set-Cookie", sessionToken)
	writer.WriteHeader(http.StatusCreated)
}
