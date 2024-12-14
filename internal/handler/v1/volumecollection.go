package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type VolumeCollectionHandler struct {
	service service.ResourceService
}

func NewVolumeCollectionHandler(service service.ResourceService) *VolumeCollectionHandler {
	return &VolumeCollectionHandler{
		service: service,
	}
}

func (handler *VolumeCollectionHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/{root:.*}/Volumes`, resourceGetter(handler.service)).Methods(http.MethodGet)

	router.HandleFunc(`/{root:.*}/Volumes`, handler.createVolume).Methods(http.MethodPost)
	router.HandleFunc(`/{root:.*}/Volumes/Members`, resourceCreatorFromNotCollectionEndpoint(handler.createVolume)).Methods(http.MethodPost)

	// Endpoint only for compatibility with official Swordfish Emulator
	//
	// Not a good way, should be removed when Ansible modules will handle
	// resources creation in right way
	router.HandleFunc(`/{root:.*}/Volumes`+idPathRegex, resourceCreatorFromNotCollectionEndpoint(handler.createVolume)).Methods(http.MethodPost)
}

func (handler *VolumeCollectionHandler) createVolume(writer http.ResponseWriter, request *http.Request) {
	volume, err := util.UnmarshalFromReader[domain.Volume](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	// Only for compatibility with official Swordfish Emulator.
	// TODO: Need to be removed when modules handle creation in right way for clearer dataflow.
	contextId := request.Context().Value(idContext)
	if contextId != nil {
		volume.Id = contextId.(string)
	}

	createdVolume, err := handler.service.AddResourceToCollection(request.Context(), dto.ResourceRequestDto{
		Name:            volume.Name,
		Id:              volume.Id,
		OdataType:       "#Volume.v1_10_0.Volume",
		Resource:        volume,
		IdSetter:        func(id string) { volume.Id = id },
		OdataIdSetter:   func(odataId string) { volume.OdataId = &odataId },
		OdataTypeSetter: func(odataType string) { volume.OdataType = &odataType },
		Collection: dto.CollectionDto{
			OdataId:   request.RequestURI,
			Name:      "Volume Collection",
			OdataType: "#VolumeCollection.VolumeCollection",
		},
	})

	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	util.WriteJSON(writer, createdVolume)
}
