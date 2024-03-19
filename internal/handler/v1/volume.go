package v1

// import (
// 	"net/http"

// 	"github.com/gorilla/mux"
// 	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
// )

// // VolumeHandler addresses to the Volume endpoints
// type VolumeHandler struct {
// 	service service.VolumeService
// }

// func NewVolumeHandler(service service.VolumeService) *VolumeHandler {
// 	return &VolumeHandler{
// 		service: service,
// 	}
// }

// // SetRouter sets handle functions for all CRUD operations on Volume resource
// func (handler *VolumeHandler) SetRouter(router *mux.Router) {
// 	router.HandleFunc("/{VolumeCollection:.*\\/Volumes}/{id:\\d+}", handler.getVolume).Methods(http.MethodGet)
// 	router.HandleFunc("/{VolumeCollection:.*\\/Volumes}/{id:\\d+}", handler.deleteVolume).Methods(http.MethodDelete)
// 	router.HandleFunc("/{VolumeCollection:.*\\/Volumes}/{id:\\d+}", handler.modifyVolume).Methods(http.MethodPatch)
// 	router.HandleFunc("/{VolumeCollection:.*\\/Volumes}/{id:\\d+}", handler.createVolume).Methods(http.MethodPost, http.MethodPut)
// }

// func (handler *VolumeHandler) getVolume(writer http.ResponseWriter, request *http.Request) {
// 	// todo
// }

// func (handler *VolumeHandler) deleteVolume(writer http.ResponseWriter, request *http.Request) {
// 	// todo
// }

// func (handler *VolumeHandler) modifyVolume(writer http.ResponseWriter, request *http.Request) {
// 	// todo
// }

// func (handler *VolumeHandler) createVolume(writer http.ResponseWriter, request *http.Request) {
// 	// todo
// }
