package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/sirupsen/logrus"
	squids "github.com/sqids/sqids-go"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
)

func Addr[T any](t T) *T {
	return &t
}

func Marshal(data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("%w.", errlib.ErrInternal)
	}

	return bytes, nil
}

func Unmarshal[T any](data []byte) (*T, error) {
	var resource T
	err := json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf("%w.", errlib.ErrInternal)
	}

	return &resource, nil
}

func UnmarshalFromReader[T any](reader io.Reader) (*T, error) {
	var resource T
	err := json.NewDecoder(reader).Decode(&resource)
	if err != nil {
		return nil, fmt.Errorf("%w.", errlib.ErrInternal)
	}

	return &resource, nil
}

func WriteJSON(writer http.ResponseWriter, jsonStruct interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(jsonStruct)
	if err != nil {
		jsonerr := errlib.GetJSONError(fmt.Errorf("%w.", errlib.ErrInternal))
		logrus.Error(err)
		writer.WriteHeader(jsonerr.Error.Code)
		if err = json.NewEncoder(writer).Encode(jsonerr); err != nil {
			logrus.Error(err)
		}
	}
}

func WriteJSONError(writer http.ResponseWriter, err error) {
	jsonErr := errlib.GetJSONError(err)
	writer.WriteHeader(jsonErr.Error.Code)
	WriteJSON(writer, jsonErr)
}

func GetParent(uri string) string {
	return filepath.ToSlash(filepath.Dir(uri))
}

func IdGenerator() func() (string, error) {
	s, _ := squids.New()
	var counter uint64 = 0
	return func() (string, error) {
		id, err := s.Encode([]uint64{counter, counter / 10, counter / 100, counter / 1000})
		counter++
		return id, err
	}
}
