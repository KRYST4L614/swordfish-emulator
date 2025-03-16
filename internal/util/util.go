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
	iface, err := net.InterfaceByName("enp0s3")
	if err != nil {
		return err
	}

	ipv4Addresses, err := initIpArressesIPv4(iface)
	if err != nil {
		return err
	}

	ethernetInterfaceType := "#EthernetInterface.v1_7_0.EthernetInterface"

	ethernetInterface := domain.EthernetInterface{
		OdataType:           &ethernetInterfaceType,
		Id:                  "1",
		Name:                iface.Name,
		IPv4Addresses:       &ipv4Addresses,
		MACAddress:          initMacAddress(iface),
		PermanentMACAddress: initPermanentMacAddress(iface),
		LinkStatus:          initEthernetInterfaceLinkStatus(),
		Status:              initEthernetInterfaceStatus(),
	}

	ethernetInterfaceJson, err := Marshal(ethernetInterface)
	if err != nil {
		return err
	}

	file, err := os.Create("datasets/basic/Systems/FileServer/EthernetInterfaces/1/index.json")
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(ethernetInterfaceJson)
	if err != nil {
		return err
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
			address, err := getIPv4Address(i)
			if err != nil {
				return nil, err
			}
			addresses = append(addresses, *address)
		}
	}
	return addresses, nil
}

func getIPv4Address(addr net.Addr) (*domain.IPAddressesV115IPv4Address, error) {
	cmd := "ip route show 0.0.0.0/0 dev enp0s3 | cut -d ' ' -f 3"

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
	subnetMask.FromIPAddressesV115SubnetMask(strings.Trim(string(mask), "\""))

	return &domain.IPAddressesV115IPv4Address{
		Address:    &ipAddr[0],
		SubnetMask: &subnetMask,
		Gateway:    &gateway,
	}, nil
}

func initMacAddress(iface *net.Interface) *domain.EthernetInterfaceV1122EthernetInterface_MACAddress {
	macAddress := domain.EthernetInterfaceV1122EthernetInterface_MACAddress{}
	macAddress.FromEthernetInterfaceV1122MACAddress(iface.HardwareAddr.String())
	return &macAddress
}

func initPermanentMacAddress(iface *net.Interface) *domain.EthernetInterfaceV1122EthernetInterface_PermanentMACAddress {
	permanentMacAddress := domain.EthernetInterfaceV1122EthernetInterface_PermanentMACAddress{}
	permanentMacAddress.FromEthernetInterfaceV1122MACAddress(iface.HardwareAddr.String())
	return &permanentMacAddress
}

func initEthernetInterfaceStatus() *domain.ResourceStatus {

	statusState := domain.ResourceStatus_State{}
	statusState.FromResourceState("Enable")

	statusHealth := domain.ResourceStatus_Health{}
	statusHealth.FromResourceHealth("OK")

	status := domain.ResourceStatus{
		State:  &statusState,
		Health: &statusHealth,
	}
	return &status
}

func initEthernetInterfaceLinkStatus() *domain.EthernetInterfaceV1122EthernetInterface_LinkStatus {
	linkStatus := domain.EthernetInterfaceV1122EthernetInterface_LinkStatus{}
	linkStatus.FromEthernetInterfaceV1122LinkStatus("LinkUp")
	return &linkStatus
}
