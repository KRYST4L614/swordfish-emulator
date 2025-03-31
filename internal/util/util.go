package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"log/slog"

	squids "github.com/sqids/sqids-go"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
)

func Addr[T any](t T) *T {
	return &t
}

func Marshal(data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("util.marshal error %w", errlib.ErrInternal)
	}

	return bytes, nil
}

func Unmarshal[T any](data []byte) (*T, error) {
	var resource T
	err := json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf("util.unmarshal error%w", errlib.ErrInternal)
	}

	return &resource, nil
}

func UnmarshalFromReader[T any](reader io.Reader) (*T, error) {
	var resource T
	err := json.NewDecoder(reader).Decode(&resource)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON %w", errlib.ErrBadRequest)
	}

	return &resource, nil
}

func WriteJSON(writer http.ResponseWriter, jsonStruct interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(jsonStruct)
	if err != nil {
		jsonerr := errlib.GetJSONError(fmt.Errorf("util.WriteJSON %w", errlib.ErrInternal))
		slog.Warn(err.Error())
		writer.WriteHeader(jsonerr.Error.Code)
		if err = json.NewEncoder(writer).Encode(jsonerr); err != nil {
			slog.Warn(err.Error())
		}
	}
}

func WriteJSONError(writer http.ResponseWriter, err error) {
	jsonErr := errlib.GetJSONError(err)
	writer.WriteHeader(jsonErr.Error.Code)
	WriteJSON(writer, jsonErr)
}

// GetParent - gets parent path
//
// Example: /foo/boo -> /foo
// works with reverse slash too: \foo\boo -> /foo
func GetParent(uri string) string {
	return filepath.ToSlash(filepath.Dir(uri))
}

// IdGenerator returns simple function that generates unique
// url-safe ids. Should be removed with service/generator.go
func IdGenerator() func() (string, error) {
	s, _ := squids.New()
	var counter uint64 = 0
	return func() (string, error) {
		id, err := s.Encode([]uint64{counter, counter / 10, counter / 100, counter / 1000})
		counter++
		return id, err
	}
}

func InitEthernetInterface() error {
	iface, err := net.Interfaces()
	if err != nil {
		return err
	}

	for _, i := range iface {
		ipv4Addresses, err := initIpArressesIPv4(&i)
		if err != nil {
			return err
		}

		ethernetInterfaceType := "#EthernetInterface.v1_7_0.EthernetInterface"

		macAddress, err := initMacAddress(&i)
		if err != nil {
			return err
		}

		permanentMacAddress, err := initPermanentMacAddress(&i)
		if err != nil {
			return err
		}

		linkStatus, err := initEthernetInterfaceLinkStatus(&i)
		if err != nil {
			return err
		}

		status, err := initEthernetInterfaceStatus(&i)
		if err != nil {
			return err
		}

		oDataId := fmt.Sprintf("/redfish/v1/Systems/FileServer/EthernetInterfaces/%v", i.Name)

		ethernetInterface := domain.EthernetInterface{
			OdataType:           &ethernetInterfaceType,
			Id:                  i.Name,
			Name:                i.Name,
			IPv4Addresses:       &ipv4Addresses,
			MACAddress:          macAddress,
			PermanentMACAddress: permanentMacAddress,
			LinkStatus:          linkStatus,
			Status:              status,
			OdataId:             &oDataId,
		}

		ethernetInterfaceJson, err := Marshal(ethernetInterface)
		if err != nil {
			return err
		}

		err = os.MkdirAll(fmt.Sprintf("datasets/basic/Systems/FileServer/EthernetInterfaces/%v/", i.Name), 0755)
		if err != nil {
			return err
		}

		ethernetInterfaceFile, err := os.Create(fmt.Sprintf("datasets/basic/Systems/FileServer/EthernetInterfaces/%v/index.json", i.Name))
		if err != nil {
			return err
		}

		defer ethernetInterfaceFile.Close()

		_, err = ethernetInterfaceFile.Write(ethernetInterfaceJson)
		if err != nil {
			return err
		}

		err = addMemberToEthernetInterfaceCollection("datasets/basic/Systems/FileServer/EthernetInterfaces/index.json", oDataId)
		if err != nil {
			return err
		}
	}

	return nil
}

func initIpArressesIPv4(iface *net.Interface) ([]domain.IPAddressesV115IPv4Address, error) {
	addresses := []domain.IPAddressesV115IPv4Address{}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	for _, i := range addrs {
		if net.ParseIP(strings.Split(i.String(), "/")[0]).To4() != nil {
			address, err := getIPv4Address(i, iface.Name)
			if err != nil {
				return nil, err
			}
			addresses = append(addresses, *address)
		}
	}
	return addresses, nil
}

func getIPv4Address(addr net.Addr, ifaceName string) (*domain.IPAddressesV115IPv4Address, error) {
	cmd := fmt.Sprintf("ip route show 0.0.0.0/0 dev %v | cut -d ' ' -f 3", ifaceName)

	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, err
	}

	gateway := strings.TrimSpace(string(out))
	ipAddr := strings.Split(addr.String(), "/")
	mask32, err := strconv.Atoi(ipAddr[1])
	if err != nil {
		return nil, err
	}

	mask, err := Marshal(net.IP(net.CIDRMask(mask32, 32)).String())
	if err != nil {
		return nil, err
	}

	subnetMask := domain.IPAddressesV115IPv4Address_SubnetMask{}
	err = subnetMask.FromIPAddressesV115SubnetMask(strings.Trim(string(mask), "\""))
	if err != nil {
		return nil, err
	}

	return &domain.IPAddressesV115IPv4Address{
		Address:    &ipAddr[0],
		SubnetMask: &subnetMask,
		Gateway:    &gateway,
	}, nil
}

func initMacAddress(iface *net.Interface) (*domain.EthernetInterfaceV1122EthernetInterface_MACAddress, error) {
	macAddress := domain.EthernetInterfaceV1122EthernetInterface_MACAddress{}
	err := macAddress.FromEthernetInterfaceV1122MACAddress(iface.HardwareAddr.String())
	if err != nil {
		return nil, err
	}
	return &macAddress, nil
}

func initPermanentMacAddress(iface *net.Interface) (*domain.EthernetInterfaceV1122EthernetInterface_PermanentMACAddress, error) {
	permanentMacAddress := domain.EthernetInterfaceV1122EthernetInterface_PermanentMACAddress{}
	err := permanentMacAddress.FromEthernetInterfaceV1122MACAddress(iface.HardwareAddr.String())
	if err != nil {
		return nil, err
	}
	return &permanentMacAddress, nil
}

func initEthernetInterfaceStatus(iface *net.Interface) (*domain.ResourceStatus, error) {
	statusState := domain.ResourceStatus_State{}
	statusHealth := domain.ResourceStatus_Health{}

	err := statusHealth.FromResourceHealth("OK")
	if err != nil {
		return nil, err
	}

	if slices.Contains(strings.Split(iface.Flags.String(), "|"), net.FlagUp.String()) {
		err = statusState.FromResourceState("Enabled")
	} else {
		err = statusState.FromResourceState("Disabled")
	}
	if err != nil {
		return nil, err
	}

	status := domain.ResourceStatus{
		State:  &statusState,
		Health: &statusHealth,
	}
	return &status, nil
}

func initEthernetInterfaceLinkStatus(iface *net.Interface) (*domain.EthernetInterfaceV1122EthernetInterface_LinkStatus, error) {
	linkStatus := domain.EthernetInterfaceV1122EthernetInterface_LinkStatus{}
	var err error
	if slices.Contains(strings.Split(iface.Flags.String(), "|"), net.FlagRunning.String()) {
		err = linkStatus.FromEthernetInterfaceV1122LinkStatus("LinkUp")
	} else {
		err = linkStatus.FromEthernetInterfaceV1122LinkStatus("LinkDown")
	}
	if err != nil {
		return nil, err
	}

	return &linkStatus, nil
}

func addMemberToEthernetInterfaceCollection(collectionPath string, oDataId string) error {
	ethernetInterfaceCollectionOldJSON, err := os.ReadFile(collectionPath)
	if err != nil {
		return err
	}

	ethernetInterfaceCollection, err := Unmarshal[domain.EthernetInterfaceCollection](ethernetInterfaceCollectionOldJSON)
	if err != nil {
		return err
	}

	if slices.ContainsFunc(*ethernetInterfaceCollection.Members, func(val domain.OdataV4IdRef) bool {
		return *(val.OdataId) == oDataId
	}) {
		return nil
	}

	*ethernetInterfaceCollection.Members = append(*ethernetInterfaceCollection.Members, domain.OdataV4IdRef{OdataId: &oDataId})
	*ethernetInterfaceCollection.MembersOdataCount = int64(len(*ethernetInterfaceCollection.Members))

	ethernetInterfaceCollectionNewJSON, err := Marshal(ethernetInterfaceCollection)
	if err != nil {
		return err
	}

	ethernetInterfaceCollectionFile, err := os.OpenFile(collectionPath, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer ethernetInterfaceCollectionFile.Close()

	_, err = ethernetInterfaceCollectionFile.Write(ethernetInterfaceCollectionNewJSON)
	if err != nil {
		return err
	}

	return nil
}
