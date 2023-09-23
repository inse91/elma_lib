package e365_gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"sync"
)

func (app App[T]) find(ctx context.Context, f filter) ([]T, error) {

	bts, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, app.method.list, bytes.NewReader(bts))
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}
	request.Header = app.stand.header()

	alr, err := doRequest[appListResponse[T]](app.client, request)
	if err != nil {
		return nil, err
	}
	if !alr.Success {
		return nil, wrap(alr.Error, ErrResponseNotSuccess)
	}

	return alr.Result.Result, nil

}

func (app App[T]) Search() searchInstance[T] {
	return searchInstance[T]{
		app:  &app,
		size: 10,
	}
}

func (s searchInstance[T]) All(ctx context.Context) ([]T, error) {
	items, err := s.app.find(ctx, filter{
		From:         s.from,
		Size:         s.size,
		Active:       !s.includeDeleted,
		SearchFilter: s.search,
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s searchInstance[T]) AllAtOnce(ctx context.Context, goroutineLimit int) ([]T, error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if goroutineLimit < 1 {
		goroutineLimit = 1
	}

	eg, _ := errgroup.WithContext(ctx)
	eg.SetLimit(goroutineLimit)

	all := make([]T, 0)
	mu := sync.Mutex{}
	for i := 0; i < 10; i++ {
		i := i
		eg.Go(func() error {

			items, err := s.app.find(ctx, filter{
				From:         i * 100,
				Size:         100,
				Active:       !s.includeDeleted,
				SearchFilter: s.search,
			})
			if err != nil {
				cancel()
				return err
			}
			if len(items) == 0 {
				//cancel()
				return ErrNoMoreItems
			}
			mu.Lock()
			defer mu.Unlock()
			all = append(all, items...)
			return nil
		})
	}

	err := eg.Wait()
	if err == nil || errors.Is(err, ErrNoMoreItems) {
		return all, nil
	}
	return nil, err
}

// First получает один элмент по переданному фильтру
func (s searchInstance[T]) First(ctx context.Context) (T, error) {
	var t T
	items, err := s.app.find(ctx, filter{
		From:         s.from,
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

type filter struct {
	From   int  `json:"from"`
	Size   int  `json:"size"`
	Active bool `json:"active"`
	SearchFilter
}

// SearchFilter - набор фильтров
type SearchFilter struct {
	Fields          Fields           `json:"filter"`
	IDs             []string         `json:"ids,omitempty"`
	SortExpressions []SortExpression `json:"sortExpressions,omitempty"`
	AtStatus        []string         `json:"statusCode,omitempty"`
	StatusGroupId   string           `json:"statusGroupId,omitempty"`
}

type searchInstance[T interface{}] struct {
	search         SearchFilter
	includeDeleted bool
	size           int
	from           int
	app            *App[T]
}

// Where применяет фильтр к поиску
func (s searchInstance[T]) Where(sf SearchFilter) searchInstance[T] {
	s.search = sf
	return s
}

func (s searchInstance[T]) Size(size int) searchInstance[T] {
	if size < 0 {
		size = 10
	}
	s.size = size
	return s
}

func (s searchInstance[T]) From(from int) searchInstance[T] {
	if from < 0 {
		from = 0
	}
	s.from = from
	return s
}

func (s searchInstance[T]) IncludeDeleted() searchInstance[T] {
	s.includeDeleted = true
	return s
}
