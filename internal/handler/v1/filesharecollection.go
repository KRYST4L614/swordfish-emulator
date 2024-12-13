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

type FileShareCollectionHandler struct {
	service service.ResourceService
}

func NewFileShareCollectionHandler(service service.ResourceService) *FileShareCollectionHandler {
	return &FileShareCollectionHandler{
		service: service,
	}
}

func (handler *FileShareCollectionHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/ExportedFileShares`, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/ExportedFileShares`, resourceGetter(handler.service)).Methods(http.MethodGet)

	router.HandleFunc(`/ExportedFileShares`, handler.createFileShare).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/ExportedFileShares`, handler.createFileShare).Methods(http.MethodPost)

	router.HandleFunc(`/ExportedFileShares`, resourceCreatorFromNotCollectionEndpoint(handler.createFileShare)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/ExportedFileShares`, resourceCreatorFromNotCollectionEndpoint(handler.createFileShare)).Methods(http.MethodPost)

	// Endpoints only for compatibility with official Swordfish Emulator
	//
	// Not a good way, should be removed when Ansible modules will handle
	// resources creation in right way
	router.HandleFunc(`/FileSystems%v/ExportedFileShares`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createFileShare)).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/FileSystems%v/ExportedFileShares`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createFileShare)).Methods(http.MethodPost)
}

func (handler *FileShareCollectionHandler) createFileShare(writer http.ResponseWriter, request *http.Request) {
	fileShare, err := util.UnmarshalFromReader[domain.FileShare](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	// Only for compatibility with official Swordfish Emulator.
	// TODO: Need to be removed when modules handle creation in right way for clearer dataflow.

	slog.Info("FileShareCollection uri: " + request.RequestURI)
	pool, err := handler.service.AddResourceToCollection(request.Context(), dto.ResourceRequestDto{
		Name:            fileShare.Name,
		Id:              fileShare.Id,
		OdataType:       "#FileShare.v1_2_0.FileShare",
		Resource:        fileShare,
		IdSetter:        func(id string) { fileShare.Id = id },
		OdataIdSetter:   func(odataId string) { fileShare.OdataId = &odataId },
		OdataTypeSetter: func(odataType string) { fileShare.OdataType = &odataType },
		Collection: dto.CollectionDto{
			OdataId:   request.RequestURI,
			Name:      "FileShareCollection",
			OdataType: "#FileShareCollection.FileShareCollection",
		},
	})

	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	util.WriteJSON(writer, pool)
}
