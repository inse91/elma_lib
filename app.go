package e365_gateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrInvalidID          = errors.New("invalid app item id")
	ErrResponseNotOK      = errors.New("response status is not ok")
	ErrResponseNotSuccess = errors.New("response is not success")
	ErrNilItem            = errors.New("item is nil")
	ErrNilSearchFilter    = errors.New("search filter is nil")
	ErrResponseNilItem    = errors.New("response item in nil")
)

const (
	methodGet       = "/get"
	methodUpdate    = "/update"
	methodCreate    = "/create"
	methodList      = "/list"
	methodSetStatus = "/set-status"
	methodGetStatus = "/settings/status"
)

type App[T interface{}] struct {
	url    string
	stand  Stand
	client *http.Client
	header http.Header
	method struct {
		create    string
		list      string
		getStatus string
	}
}

// NewApp creates new adapter for interaction with app in elma, where T is app Context
func NewApp[T interface{}](settings AppSettings) App[T] {
	url := settings.toAppUrl()
	return App[T]{
		stand: settings.Stand,
		url:   url,
		client: &http.Client{
			Timeout: time.Second * 5,
		},
		method: struct {
			create    string
			list      string
			getStatus string
		}{
			create:    fmt.Sprintf("%s%s", url, methodCreate),
			list:      fmt.Sprintf("%s%s", url, methodList),
			getStatus: fmt.Sprintf("%s%s", url, methodGetStatus),
		},
		//header: func() http.Header {
		//	h := http.Header{}
		//	h.Set("Content-type", "application/json")
		//	h.Set("Authorization", settings.Stand.token())
		//	return h
		//}(),
	}
}

// Update updates app item in elma by given __id
func (app App[T]) Update(id string, item T) (T, error) {
	var nilT T
	//if item == nilT {
	//	return nilT, ErrNilItem
	//}

	if len(id) != uuid4Len {
		return nilT, fmt.Errorf("%s: %w", id, ErrInvalidID)
	}

	url := app.url + "/" + id + methodUpdate
	bts, err := json.Marshal(createItemRequest[T]{
		Context: item,
	})
	if err != nil {
		return nilT, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bts))
	if err != nil {
		return nilT, fmt.Errorf("failed creating request: %w", err)
	}
	request.Header = app.stand.header()

	response, err := app.client.Do(request)
	if err != nil {
		return nilT, fmt.Errorf("failed sending request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	ir := new(itemResponse[T])
	if err = decodeStd(response.Body, ir); err != nil {
		return nilT, fmt.Errorf("failed decoding response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nilT, fmt.Errorf("%w: %d (%s)", ErrResponseNotOK, response.StatusCode, ir.Error)
	}

	if !ir.Success {
		return nilT, fmt.Errorf("%w: %s", ErrResponseNotSuccess, ir.Error)
	}

	//if ir.Item == nil {
	//	return nilT, ErrResponseNilItem
	//}
	return ir.Item, nil
}

// Create creates app item in elma
func (app App[T]) Create(item T) (T, error) {

	var nilT T

	//if item == nil {
	//	return nilT, ErrNilItem
	//}

	bts, err := json.Marshal(createItemRequest[T]{
		Context: item,
	})
	if err != nil {
		return nilT, err
	}

	request, err := http.NewRequest(http.MethodPost, app.method.create, bytes.NewReader(bts))
	if err != nil {
		return nilT, fmt.Errorf("failed creating request: %w", err)
	}

	request.Header = app.stand.header()
	response, err := app.client.Do(request)
	if err != nil {
		return nilT, fmt.Errorf("failed sending request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	ir := new(itemResponse[T])
	if err = decodeStd(response.Body, ir); err != nil {
		return nilT, fmt.Errorf("failed decoding response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nilT, fmt.Errorf("%w: %d (%s)", ErrResponseNotOK, response.StatusCode, ir.Error)
	}

	if !ir.Success {
		return nilT, fmt.Errorf("%w: %s", ErrResponseNotSuccess, ir.Error)
	}

	//if ir.Item == nil {
	//	return nilT, ErrResponseNilItem
	//}

	return ir.Item, nil

}

// GetByID performs search app item by given __id
func (app App[T]) GetByID(id string) (T, error) {
	var nilT T
	if len(id) != uuid4Len {
		return nilT, fmt.Errorf("%s: %w", id, ErrInvalidID)
	}

	url := app.url + "/" + id + methodGet
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nilT, fmt.Errorf("failed creating request: %w", err)
	}
	request.Header = app.header

	response, err := app.client.Do(request)
	if err != nil {
		return nilT, fmt.Errorf("failed sending request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	ir := new(itemResponse[T])
	if err = decodeStd(response.Body, ir); err != nil {
		return nilT, fmt.Errorf("failed decoding response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nilT, fmt.Errorf("%w: %d (%s)", ErrResponseNotOK, response.StatusCode, ir.Error)
	}

	if !ir.Success {
		return nilT, fmt.Errorf("%w: %s", ErrResponseNotSuccess, ir.Error)
	}

	//if ir.Item == nil {
	//	return nilT, ErrResponseNilItem
	//}

	return ir.Item, nil

}

// SetStatus sets app item status with __id by given status code
func (app App[T]) SetStatus(id, code string) (T, error) {

	var nilT T
	if len(id) != uuid4Len {
		return nilT, fmt.Errorf("%s: %w", id, ErrInvalidID)
	}

	url := app.url + "/" + id + methodSetStatus
	bts, err := json.Marshal(setStatusRequest{
		Status: statusCode{
			Code: code,
		},
	})
	if err != nil {
		return nilT, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bts))
	if err != nil {
		return nilT, fmt.Errorf("failed creating request: %w", err)
	}
	request.Header = app.header

	response, err := app.client.Do(request)
	if err != nil {
		return nilT, fmt.Errorf("failed sending request: %w", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	ir := new(itemResponse[T])
	if err = decodeStd(response.Body, ir); err != nil {
		return nilT, fmt.Errorf("failed decoding response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nilT, fmt.Errorf("%w: %d (%s)", ErrResponseNotOK, response.StatusCode, ir.Error)
	}

	if !ir.Success {
		return nilT, fmt.Errorf("%w: %s", ErrResponseNotSuccess, ir.Error)
	}

	//if ir.Item == nil {
	//	return nilT, ErrResponseNilItem
	//}

	return ir.Item, nil
}

// GetStatusInfo gets app status variants
func (app App[T]) GetStatusInfo() (StatusInfo, error) {

	request, err := http.NewRequest(http.MethodGet, app.method.getStatus, nil)
	if err != nil {
		return StatusInfo{}, fmt.Errorf("failed creating request: %w", err)
	}
	request.Header = app.header

	response, err := app.client.Do(request)
	if err != nil {
		return StatusInfo{}, fmt.Errorf("failed sending request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	gsr := new(getStatusResponse)
	if err = decodeStd(response.Body, gsr); err != nil {
		return StatusInfo{}, fmt.Errorf("failed decoding response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return StatusInfo{}, fmt.Errorf("%w: %d (%s)", ErrResponseNotOK, response.StatusCode, gsr.Error)
	}

	if !gsr.Success {
		return StatusInfo{}, fmt.Errorf("%w: %s", ErrResponseNotSuccess, gsr.Error)
	}

	return gsr.StatusInfo, nil

}
