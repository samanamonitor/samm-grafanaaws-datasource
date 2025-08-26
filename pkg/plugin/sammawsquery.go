package plugin

import (
    "errors"

    "github.com/grafana/grafana-plugin-sdk-go/backend"
)

type SammAwsQuery interface {
    ListActions()   ([]byte, error)
    CallAction()    ([]byte, error)
    QueryVariable() ([]byte, error)
    QueryData()     backend.DataResponse
}

type NotImplemented struct {}

func (NotImplemented) QueryData() backend.DataResponse {
    return backend.ErrDataResponse(backend.StatusBadRequest, "Not implemented")
}

func (NotImplemented) QueryVariable() ([]byte, error) {
    return []byte{}, errors.New("Not Implemented")
}

func (NotImplemented) ListActions() ([]byte, error) {
    return []byte{}, errors.New("Not Implemented")
}

func (NotImplemented) CallAction() ([]byte, error) {
    return []byte{}, errors.New("Not Implemented")
}

type SammAwsResponse struct {
    Message string `json:"message"`
}

