package util

import (
	"encoding/json"
	"fmt"

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
