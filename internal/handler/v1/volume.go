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

type VolumeHandler struct {
	service service.ResourceService
}

func NewVolumeHandler(service service.ResourceService) *VolumeHandler {
	return &VolumeHandler{
		service: service,
	}
}

func (handler *VolumeHandler) SetRouter(router *mux.Router) {
	router.HandleFunc(`/{root:.*}/Volumes`+idPathRegex, resourceGetter(handler.service)).Methods(http.MethodGet)
	router.HandleFunc(`/{root:.*}/Volumes`+idPathRegex, handler.updateVolume).Methods(http.MethodPatch)
	router.HandleFunc(`/{root:.*}/Volumes`+idPathRegex, handler.replaceVolume).Methods(http.MethodPut)
	router.HandleFunc(`/{root:.*}/Volumes`+idPathRegex, resourceDeleter(handler.service)).Methods(http.MethodDelete)
}

func (handler *VolumeHandler) replaceVolume(writer http.ResponseWriter, request *http.Request) {
	volume, err := util.UnmarshalFromReader[domain.Volume](request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	volumeId := request.RequestURI
	volume.Id = filepath.Base(volumeId)
	*volume.OdataId = volumeId

	newVolume, err := handler.service.Replace(request.Context(), volumeId, volume)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newVolume)
}

func (handler *VolumeHandler) updateVolume(writer http.ResponseWriter, request *http.Request) {
	patchData, err := io.ReadAll(request.Body)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	volumeId := request.RequestURI
	newVolume, err := handler.service.Update(request.Context(), volumeId, patchData)
	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	util.WriteJSON(writer, newVolume)
}
