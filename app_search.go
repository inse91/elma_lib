package e365_gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (app App[T]) find(f filter) ([]T, error) {

	bts, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, app.method.list, bytes.NewReader(bts))
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}
	request.Header = app.stand.header()

	response, err := app.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed sending request: %w", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	//bts1, err := io.ReadAll(response.Body)
	//_ = bts1

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrResponseNotOK, response.StatusCode)
	}

	alr := new(appListResponse[T])
	if err = decodeStd(response.Body, alr); err != nil {
		return nil, fmt.Errorf("failed decoding response body: %w", err)
	}

	if !alr.Success {
		return nil, fmt.Errorf("%w: %s", ErrResponseNotSuccess, alr.Error)
	}

	return alr.Result.Result, nil

}

func (app App[T]) Search(sf ...SearchFilter) searchInstance[T] {
	return searchInstance[T]{
		search: func() SearchFilter {
			if len(sf) == 0 {
				return SearchFilter{}
			}
			return sf[0]
		}(),
		app: &app,
	}
}

func (s searchInstance[T]) All() ([]T, error) {
	items, err := s.app.find(filter{
		From:         0,
		Size:         100,
		Active:       !s.includeDeleted,
		SearchFilter: s.search,
	})
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s searchInstance[T]) First() (T, error) {
	var t T
	items, err := s.app.find(filter{
		From:         0,
		Size:         1,
		Active:       !s.includeDeleted,
		SearchFilter: s.search,
	})
	if err != nil {
		return t, err
	}

	if len(items) == 0 {
		return t, nil
	}

	return items[0], nil
}

func AppDateFilter(min, max string) map[string]string {
	return map[string]string{
		"min": min,
		"max": max,
	}
}

func AppNumberFilter(min, max float64) map[string]float64 {
	return map[string]float64{
		"min": min,
		"max": max,
	}
}

type filter struct {
	From   int  `json:"from"`
	Size   int  `json:"size"`
	Active bool `json:"active"`
	SearchFilter
}

type SearchFilter struct {
	Fields          Fields           `json:"filter"`
	IDs             []string         `json:"ids,omitempty"`
	SortExpressions []SortExpression `json:"sortExpressions,omitempty"`
	InStatuses      []string         `json:"statusCode,omitempty"`
	StatusGroupId   string           `json:"statusGroupId,omitempty"`
}

type searchInstance[T interface{}] struct {
	search         SearchFilter
	includeDeleted bool
	app            *App[T]
}

func (s searchInstance[T]) IncludeDeleted() searchInstance[T] {
	//s.filter.Active = false
	s.includeDeleted = true
	return s
}

type SortExpression struct {
	Ascending bool   `json:"ascending"`
	Field     string `json:"field"`
}

type Fields map[string]interface{}

func (f Fields) MarshalJSON() ([]byte, error) {
	const emptyTf = "{\"tf\":{}}"
	l := len(f)
	if f == nil || l == 0 {
		return []byte(emptyTf), nil
	}

	sb := new(strings.Builder)
	sb.WriteString("{\"tf\":{")
	var bts []byte
	var err error
	i := 1
	for k, v := range f {
		sb.WriteRune('"')
		sb.WriteString(k)
		sb.WriteRune('"')
		sb.WriteRune(':')
		if bts, err = json.Marshal(v); err != nil {
			return nil, err
		}
		sb.Write(bts)
		if i != l {
			sb.WriteRune(',')
		}
	}
	sb.WriteString("}}")
	return []byte(sb.String()), nil
}
