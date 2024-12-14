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

type FileSystemHandler struct {
	service service.ResourceService
}

func NewFileSystemHandler(service service.ResourceService) *FileSystemHandler {
	return &FileSystemHandler{
		service: service,
	}
}

func (handler *FileSystemHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/FileSystems`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/FileSystems`+idPathRegex, handler.updateFileSystem).Methods(http.MethodPatch)
	router.HandleFunc(`/FileSystems`+idPathRegex, handler.replaceFileSystem).Methods(http.MethodPut)
	router.HandleFunc(`/FileSystems`+idPathRegex, resourceDeleter(handler.service)).Methods(http.MethodDelete)

	router.HandleFunc(`/{root:.*}/FileSystems`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/FileSystems`+idPathRegex, handler.updateFileSystem).Methods(http.MethodPatch)
	router.HandleFunc(`/{root:.*}/FileSystems`+idPathRegex, handler.replaceFileSystem).Methods(http.MethodPut)
	router.HandleFunc(`/{root:.*}/FileSystem`+idPathRegex, resourceDeleter(handler.service)).Methods(http.MethodDelete)
}

func (handler *FileSystemHandler) replaceFileSystem(writer http.ResponseWriter, request *http.Request) {
	fileSystem, err := util.UnmarshalFromReader[domain.FileSystem](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	fileSystemId := request.RequestURI
	fileSystem.Id = filepath.Base(fileSystemId)
	*fileSystem.OdataId = fileSystemId

	newFileSystem, err := handler.service.Replace(request.Context(), fileSystemId, fileSystem)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newFileSystem)
}

func (handler *FileSystemHandler) updateFileSystem(writer http.ResponseWriter, request *http.Request) {
	patchData, err := io.ReadAll(request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	fileSystemId := request.RequestURI
	newFileSystem, err := handler.service.Update(request.Context(), fileSystemId, patchData)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newFileSystem)
}
