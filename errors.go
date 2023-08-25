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

	ErrInvalidID          = errors.New("invalid item id")
	ErrResponseNotSuccess = errors.New("response is not success")
	ErrNilItem            = errors.New("item is nil")
	ErrEmptyBuffer        = errors.New("item is nil")
	ErrNilSearchFilter    = errors.New("search filter is nil")
	ErrResponseNilItem    = errors.New("response item in nil")

	ErrCreateFormData       = errors.New("failed creating form data")
	ErrWriteBytesBuffer     = errors.New("failed writing from bytes buffer to form data")
	ErrCloseMultipartWriter = errors.New("failed closing multipart writer")
)

func wrap(msg string, err error) error {
	return fmt.Errorf("%w: %s", err, msg)
}
