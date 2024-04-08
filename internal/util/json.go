package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
)

func Marshal(data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("%w. %s", errlib.ErrInternal, err.Error())
	}

	return bytes, nil
}

func Unmarshal[T any](data []byte) (*T, error) {
	var resource T
	err := json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf("%w. %s", errlib.ErrInternal, err.Error())
	}

	return &resource, nil
}

func WriteJSON(writer http.ResponseWriter, jsonStruct interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(jsonStruct)
	if err != nil {
		jsonerr := errlib.GetJSONError(fmt.Errorf("%w. %s", errlib.ErrInternal, err.Error()))
		logrus.Error(err)
		writer.WriteHeader(jsonerr.Error.Code)

		if err = json.NewEncoder(writer).Encode(jsonerr); err != nil {
			logrus.Error(err)
		}
	}
}
