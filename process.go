package e365_gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	methodRun = "/run"
)

type EmptyCtx struct{}

type Process[T interface{}] struct {
	url    string
	stand  Stand
	client *http.Client
	method struct {
		run string
	}
}

// NewProcess creates new adapter for interaction with process in elma, where T is process Context
func NewProcess[T interface{}](settings AppSettings) Process[T] {
	return Process[T]{
		url: settings.toBpmUrl(),
		client: &http.Client{
			Timeout: time.Second * 3,
		},
		stand:  settings.Stand,
		method: struct{ run string }{run: settings.toBpmUrl() + methodRun},
	}
}

func (proc Process[T]) Run(ctx T) (string, error) {

	bts, err := json.Marshal(runProcRequest[T]{
		Context: ctx,
	})
	if err != nil {
		return "", fmt.Errorf("failed encoding request body: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, proc.method.run, bytes.NewReader(bts))
	if err != nil {
		return "", fmt.Errorf("failed creating request: %w", err)
	}

	request.Header = proc.stand.header()
	response, err := proc.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed sending request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: %d", ErrResponseNotOK, response.StatusCode)
	}

	ir := new(runProcessResponse[T])
	if err = decodeStd(response.Body, ir); err != nil {
		return "", fmt.Errorf("failed decoding response body: %w", err)
	}

	if !ir.Success {
		return "", fmt.Errorf("%w: %s", ErrResponseNotSuccess, ir.Error)
	}

	return ir.Context.Id, nil

}
