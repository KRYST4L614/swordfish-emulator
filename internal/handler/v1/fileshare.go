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

type FileShareHandler struct {
	service service.ResourceService
}

func NewFileShareHandler(service service.ResourceService) *FileShareHandler {
	return &FileShareHandler{
		service: service,
	}
}

func (handler *FileShareHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/ExportedFileShares`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/ExportedFileShares`+idPathRegex, handler.updateFileShare).Methods(http.MethodPatch)
	router.HandleFunc(`/ExportedFileShares`+idPathRegex, handler.replaceFileShare).Methods(http.MethodPut)
	router.HandleFunc(`/ExportedFileShares`+idPathRegex, resourceDeleter(handler.service)).Methods(http.MethodDelete)

	router.HandleFunc(`/{root:.*}/ExportedFileShares`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/ExportedFileShares`+idPathRegex, handler.updateFileShare).Methods(http.MethodPatch)
	router.HandleFunc(`/{root:.*}/ExportedFileShares`+idPathRegex, handler.replaceFileShare).Methods(http.MethodPut)
	router.HandleFunc(`/{root:.*}/ExportedFileShares`+idPathRegex, resourceDeleter(handler.service)).Methods(http.MethodDelete)
}

func (handler *FileShareHandler) replaceFileShare(writer http.ResponseWriter, request *http.Request) {
	fileShare, err := util.UnmarshalFromReader[domain.FileShare](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	fileShareId := request.RequestURI
	fileShare.Id = filepath.Base(fileShareId)
	*fileShare.OdataId = fileShareId

	newFileShare, err := handler.service.Replace(request.Context(), fileShareId, fileShare)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newFileShare)
}

func (handler *FileShareHandler) updateFileShare(writer http.ResponseWriter, request *http.Request) {
	patchData, err := io.ReadAll(request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	fileShareId := request.RequestURI
	newFileShare, err := handler.service.Update(request.Context(), fileShareId, patchData)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newFileShare)
}
