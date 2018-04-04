package apiv1

import (
	"encoding/json"
	"errors"
)

type CloudPropsImpl struct {
	json.RawMessage
}

var _ json.Marshaler = CloudPropsImpl{}

func (p CloudPropsImpl) As(val interface{}) error {
	return json.Unmarshal([]byte(p.RawMessage), val)
}

func (c CloudPropsImpl) MarshalJSON() ([]byte, error) {
	return nil, errors.New("Expected to not marshal CloudPropsImpl as JSON")
}

func (c CloudPropsImpl) _final() {}
