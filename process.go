package e365_gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const (
	pubV1ApiInstance = "pub/v1/bpm/instance/"
	pubV1ApiBpm      = "pub/v1/bpm/template"
	methodRun        = "/run"
)

// EmptyProcCtx - пустой контекст бизнес-процесса, содержит только служебные поля.
type EmptyProcCtx struct {
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

// NewProc создает новый адаптер для взаимодействия с бизнес-процессом.
// Параметр T предстваляет собой контекст процесса. Служебные поля контекста можно взять из ProcCommon
// и встроить в структуру с пользовательскими полями контекста.
// При пстуом входном контексте можно использовать EmptyProcCtx.
// Параметры процесса передаются через Settings:
// Stand - интерфейс стенда, на котором нужно запускать процесс (!= nil);
// Namespace - код раздела/приложения, в котором находится процесс (если процесс находится на уровне раздела X,
// то код нужно передавать как  "X"; если процесс находится на уровне приложения Y в разделе X, то код нужно передавать
// как "X.Y";
// Code - код самого процесса
func NewProc[T interface{}](settings Settings) Proc[T] {
	return Proc[T]{
		url: settings.toBpmUrl(),
		client: &http.Client{
			Timeout: time.Second * 3,
		},
		stand:  settings.Stand,
		method: struct{ run string }{run: settings.toBpmUrl() + methodRun},
	}
}

// SetClientTimeout устанавливает таймаут ожидания ответа на запрос.
func (proc Proc[T]) SetClientTimeout(t time.Duration) {
	proc.client.Timeout = t
}

// GetInstanceById получает экзмемпляр бизнес-процесса по __id.
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

	ir, err := doRequest[getProcInstanceResponse[T]](proc.client, request)
	if err != nil {
		return nilT, err
	}
	if !ir.Success {
		return nilT, wrap(ir.Error, ErrResponseNotSuccess)
	}

	return ir.Context, nil

}

// Run запускает бизнес-процесс с переданным входным контекстом.
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

	ir, err := doRequest[runProcResponse[T]](proc.client, request)
	if err != nil {
		return nilT, err
	}
	if !ir.Success {
		return nilT, wrap(ir.Error, ErrResponseNotSuccess)
	}

	return ir.Context, nil

}
