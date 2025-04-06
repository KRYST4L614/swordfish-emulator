package v1

import (
	"net/http"

	"log/slog"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type FileSystemCollectionHandler struct {
	service service.ResourceService
}

func NewFileSystemCollectionHandler(service service.ResourceService) *FileSystemCollectionHandler {
	return &FileSystemCollectionHandler{
		service: service,
	}
}

func (handler *FileSystemCollectionHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/{root:.*}/FileSystems`, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/FileSystems`, resourceGetter(handler.service)).Methods(http.MethodGet)

	router.HandleFunc(`/FileSystems`, handler.createFileSystem).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/FileSystems`, handler.createFileSystem).Methods(http.MethodPost)

	router.HandleFunc(`/FileSystems`, resourceCreatorFromNotCollectionEndpoint(handler.createFileSystem)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/FileSystems`, resourceCreatorFromNotCollectionEndpoint(handler.createFileSystem)).Methods(http.MethodPost)

	// Endpoints only for compatibility with official Swordfish Emulator
	//
	// Not a good way, should be removed when Ansible modules will handle
	// resources creation in right way
	router.HandleFunc(`/FileSystems`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createFileSystem)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/FileSystems`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createFileSystem)).Methods(http.MethodPost)
}

func (handler *FileSystemCollectionHandler) createFileSystem(writer http.ResponseWriter, request *http.Request) {
	fileSystem, err := util.UnmarshalFromReader[domain.FileSystem](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	// Only for compatibility with official Swordfish Emulator.
	// TODO: Need to be removed when modules handle creation in right way for clearer dataflow.

	slog.Info("FileSystemCollection uri: " + request.RequestURI)
	createdFileSystem, err := handler.service.AddResourceToCollection(request.Context(), dto.ResourceRequestDto{
		Name:            fileSystem.Name,
		Id:              fileSystem.Name,
		OdataType:       "#FileSystem.v1_4_0.FileSystem",
		Resource:        fileSystem,
		IdSetter:        func(id string) { fileSystem.Id = id },
		OdataIdSetter:   func(odataId string) { fileSystem.OdataId = &odataId },
		OdataTypeSetter: func(odataType string) { fileSystem.OdataType = &odataType },
		Collection: dto.CollectionDto{
			OdataId:   request.RequestURI,
			Name:      "FileSystem",
			OdataType: "#FileSystemCollection.FileSystemCollection",
		},
	})

	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	util.WriteJSON(writer, createdFileSystem)
}
