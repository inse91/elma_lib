package e365_gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	pubV1ApiApp = "pub/v1/app"

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
func NewApp[T interface{}](settings Settings) App[T] {
	url := settings.toAppUrl()
	return App[T]{
		stand: settings.Stand,
		url:   url,
		client: &http.Client{
			Timeout: time.Second * 5,
			//Timeout: time.Millisecond * 100,
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
	}
}

// Create creates app item in elma
func (app App[T]) Create(ctx context.Context, item T) (T, error) {

	var nilT T
	bts, err := json.Marshal(createItemRequest[T]{
		Context: item,
	})
	if err != nil {
		return nilT, wrap(err.Error(), ErrEncodeRequestBody)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, app.method.create, bytes.NewReader(bts))
	if err != nil {
		return nilT, wrap(err.Error(), ErrCreateRequest)
	}
	request.Header = app.stand.header()

	ir, err := doRequest[itemResponse[T]](app.client, request)
	if err != nil {
		return nilT, err
	}

	if !ir.Success {
		return nilT, wrap(ir.Error, ErrResponseNotSuccess)
	}

	return ir.Item, nil

}

// GetByID получает экземпляр приложения с переданным id
func (app App[T]) GetByID(ctx context.Context, id string) (T, error) {
	var nilT T
	if len(id) != uuid4Len {
		return nilT, wrap(id, ErrInvalidID)
	}

	url := app.url + "/" + id + methodGet
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		if err != nil {
			return nilT, wrap(err.Error(), ErrCreateRequest)
		}
	}
	request.Header = app.stand.header()

	ir, err := doRequest[itemResponse[T]](app.client, request)
	if err != nil {
		return nilT, err
	}

	if !ir.Success {
		return nilT, wrap(ir.Error, ErrResponseNotSuccess)
	}
	return ir.Item, nil

}

// Update обновляет экземпляр приложения с переданным id
func (app App[T]) Update(ctx context.Context, id string, item T) (T, error) {

	var nilT T
	if len(id) != uuid4Len {
		return nilT, wrap(id, ErrInvalidID)
	}

	url := app.url + "/" + id + methodUpdate
	bts, err := json.Marshal(createItemRequest[T]{
		Context: item,
	})
	if err != nil {
		return nilT, wrap(err.Error(), ErrEncodeRequestBody)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bts))
	if err != nil {
		return nilT, wrap(err.Error(), ErrCreateRequest)
	}
	request.Header = app.stand.header()

	ir, err := doRequest[itemResponse[T]](app.client, request)
	if err != nil {
		return nilT, err
	}

	if !ir.Success {
		return nilT, wrap(ir.Error, ErrResponseNotSuccess)
	}

	return ir.Item, nil
}

// SetStatus меняет статус экземпляр приложения с переданным id на статус code
func (app App[T]) SetStatus(ctx context.Context, id, code string) (T, error) {

	var nilT T
	if len(id) != uuid4Len {
		return nilT, wrap(id, ErrInvalidID)
	}

	url := app.url + "/" + id + methodSetStatus
	bts, err := json.Marshal(setStatusRequest{
		Status: statusCode{
			Code: code,
		},
	})
	if err != nil {
		return nilT, wrap(err.Error(), ErrEncodeRequestBody)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bts))
	if err != nil {
		return nilT, wrap(err.Error(), ErrCreateRequest)
	}
	request.Header = app.stand.header()

	ir, err := doRequest[itemResponse[T]](app.client, request)
	if err != nil {
		return nilT, err
	}

	if !ir.Success {
		return nilT, wrap(ir.Error, ErrResponseNotSuccess)
	}

	return ir.Item, nil
}

// GetStatusInfo получает информацию о возможных статусах приложения
func (app App[T]) GetStatusInfo(ctx context.Context) (StatusInfo, error) {

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, app.method.getStatus, nil)
	if err != nil {
		return StatusInfo{}, wrap(err.Error(), ErrCreateRequest)
	}
	request.Header = app.stand.header()

	gsr, err := doRequest[getStatusResponse](app.client, request)
	if err != nil {
		return StatusInfo{}, err
	}

	if !gsr.Success {
		return StatusInfo{}, wrap(gsr.Error, ErrResponseNotSuccess)
	}

	return gsr.StatusInfo, nil

}
