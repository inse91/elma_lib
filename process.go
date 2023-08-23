package e365_gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const (
	methodRun = "/run"
	//methodGet = "/run"
)

type EmptyCtx struct {
	ProcCommon
}

type Proc[T interface{}] struct {
	url    string
	stand  Stand
	client *http.Client
	method struct {
		run string
	}
}

// NewProc creates new adapter for interaction with process in elma, where T is process Context
func NewProc[T interface{}](settings AppSettings) Proc[T] {
	return Proc[T]{
		url: settings.toBpmUrl(),
		client: &http.Client{
			Timeout: time.Second * 3,
		},
		stand:  settings.Stand,
		method: struct{ run string }{run: settings.toBpmUrl() + methodRun},
	}
}

func (proc Proc[T]) GetInstanceById(ctx context.Context, id string) (T, error) {

	var nilT T

	if len(id) != uuid4Len {
		return nilT, wrap(id, ErrInvalidID)
	}

	url := proc.stand.url() + "/" + pubV1ApiInstance + id + methodGet
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nilT, wrap(err.Error(), ErrCreateRequest)
	}

	request.Header = proc.stand.header()
	response, err := proc.client.Do(request)
	if err != nil {
		return nilT, wrap(err.Error(), ErrSendRequest)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nilT, wrap(response.Status, ErrResponseStatusNotOK)
	}

	ir := new(getProcInstanceResponse[T])
	if err = decodeStd(response.Body, ir); err != nil {
		return nilT, wrap(err.Error(), ErrDecodeResponseBody)
	}

	if !ir.Success {
		return nilT, wrap(ir.Error, ErrResponseNotSuccess)
	}

	return ir.Context, nil

}

func (proc Proc[T]) Run(ctx context.Context, procCtx T) (T, error) {

	var nilT T
	bts, err := json.Marshal(runProcRequest[T]{
		Context: procCtx,
	})
	if err != nil {
		return nilT, wrap(err.Error(), ErrEncodeRequestBody)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, proc.method.run, bytes.NewReader(bts))
	if err != nil {
		return nilT, wrap(err.Error(), ErrCreateRequest)
	}

	request.Header = proc.stand.header()
	response, err := proc.client.Do(request)
	if err != nil {
		return nilT, wrap(err.Error(), ErrSendRequest)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nilT, wrap(err.Error(), ErrResponseStatusNotOK)
	}

	ir := new(runProcResponse[T])
	if err = decodeStd(response.Body, ir); err != nil {
		return nilT, wrap(err.Error(), ErrDecodeResponseBody)
	}

	if !ir.Success {
		return nilT, wrap(ir.Error, ErrResponseNotSuccess)
	}

	return ir.Context, nil

}
