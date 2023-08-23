package e365_gateway

import (
	"errors"
	"fmt"
)

var (
	ErrCreateRequest       = errors.New("failed creating new http request")
	ErrSendRequest         = errors.New("failed sending http request")
	ErrResponseStatusNotOK = errors.New("response status is not ok")
	ErrDecodeResponseBody  = errors.New("failed decoding response body")
	ErrEncodeRequestBody   = errors.New("failed encoding request body")

	ErrInvalidID          = errors.New("invalid app item id")
	ErrResponseNotSuccess = errors.New("response is not success")
	ErrNilItem            = errors.New("item is nil")
	ErrNilSearchFilter    = errors.New("search filter is nil")
	ErrResponseNilItem    = errors.New("response item in nil")
)

func wrap(msg string, err error) error {
	return fmt.Errorf("%w: %s", err, msg)
}
