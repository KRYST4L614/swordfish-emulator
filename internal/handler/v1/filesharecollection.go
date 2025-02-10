package v1

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

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

	ethernetInterface, err := handler.service.Get(request.Context(), *fileShare.EthernetInterfaces.OdataId)

	if err != nil {
		util.WriteJSONError(writer, err)
		return
	}

	a, _ := util.Marshal(ethernetInterface)

	ethernetInterfaceDto, ok := util.Unmarshal[domain.EthernetInterface](a)
	if ok != nil {
		util.WriteJSONError(writer, fmt.Errorf("type cast error"))
		return
	}

	for _, ip := range *ethernetInterfaceDto.IPv4Addresses {
		subnet, _ := ip.SubnetMask.AsIPAddressesV115SubnetMask()
		mask := net.IPMask(net.ParseIP(subnet).To4())
		prefixSize, _ := mask.Size()
		if err := mountFS(*fileShare, fmt.Sprintf("%s/%v", *ip.Address, prefixSize)); err != nil {
			util.WriteJSONError(writer, err)
			return
		}
	}

	createdFileShare, err := handler.service.AddResourceToCollection(request.Context(), dto.ResourceRequestDto{
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
	util.WriteJSON(writer, createdFileShare)
}

// mount FileShare
func mountFS(fileShare domain.FileShare, ip string) error {
	isNew, err := createDirectory(*fileShare.FileSharePath)
	if err != nil {
		return err
	}

	if err = configureNFSExport(*fileShare.FileSharePath, ip); err != nil && isNew {
		if removeError := removeDirectory(*fileShare.FileSharePath); removeError != nil {
			return removeError
		}
		return err
	}

	if err != nil {
		return err
	}

	if err := exportFS(); err != nil && isNew {
		if removeError := removeDirectory(*fileShare.FileSharePath); removeError != nil {
			return removeError
		}
		if clearError := clearNFSExport(); clearError != nil {
			return clearError
		}
		return err
	}

	if err != nil {
		if clearError := clearNFSExport(); clearError != nil {
			return clearError
		}
		return err
	}

	return nil
}

// create Directory if not exists
func createDirectory(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		slog.Info(fmt.Sprintf("Creating directory: %s\n", path))
		return true, os.MkdirAll(path, 0777)
	}
	return false, nil
}

// remove Directory if not exists
func removeDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("directory doesn't exists")
	}

	return os.RemoveAll(path)
}

// configureNFSExport add note about FileShare to /etc/exports
func configureNFSExport(sharePath, clientIP string) error {
	exportEntry := fmt.Sprintf("%s %s(rw,sync,no_subtree_check)\n", sharePath, clientIP)
	exportsFile := "/etc/exports"

	//Открываем файл /etc/exports для добавления записи
	file, err := os.OpenFile(exportsFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("failed to open file /etc/exports: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(exportEntry); err != nil {
		return fmt.Errorf("failed to write to file /etc/exports: %v", err)
	}

	slog.Info(fmt.Sprintf("Note about FileShare added to /etc/export: %v", exportEntry))
	return nil
}

// clearNFSExport remove last line from /etc/exports
func clearNFSExport() error {
	exportsFile := "/etc/exports"

	//Открываем файл /etc/exports для добавления записи
	content, err := os.ReadFile(exportsFile)
	if err != nil {
		return fmt.Errorf("file read  %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) > 0 {
		lines = lines[:len(lines)-1]
	}

	newContent := strings.Join(lines, "\n") + "\n"

	if err := os.WriteFile(exportsFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("file write error %w", err)
	}

	return nil
}

// Export FS
func exportFS() error {
	slog.Info("Exporting FS...")
	cmd := exec.Command("exportfs", "-ra")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("exportfs error: %v", err)
	}
	return nil
}
