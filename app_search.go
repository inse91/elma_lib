package e365_gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"sync"
)

// find приватный общий метод для поиска, используется под капотом во всех публичных методах: First, All, AllAtOnce
func (app App[T]) find(ctx context.Context, f filter) ([]T, int, error) {

	bts, err := json.Marshal(f)
	if err != nil {
		return nil, 0, wrap(err.Error(), ErrEncodeRequestBody)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, app.method.list, bytes.NewReader(bts))
	if err != nil {
		return nil, 0, wrap(err.Error(), ErrCreateRequest)
	}
	request.Header = app.stand.header()

	alr, err := doRequest[appListResponse[T]](app.client, request)
	if err != nil {
		return nil, 0, err
	}
	if !alr.Success {
		return nil, 0, wrap(alr.Error, ErrResponseNotSuccess)
	}

	return alr.Result.Result, alr.Result.Total, nil

}

// Search используется для вызова конструктора поиска
func (app App[T]) Search() searchInstance[T] {
	return searchInstance[T]{
		app:  &app,
		size: 10,
	}
}

// All получает все элементы по переданному фильтру.
// Кол-во элементов задается через Size (не более 100, по умолчанию 10)
func (s searchInstance[T]) All(ctx context.Context) ([]T, error) {
	items, _, err := s.app.find(ctx, filter{
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

// AllAtOnce получает все элементы по переданному фильтру, кол-во элементов задается через Size (в том числе более 100).
// Асинхронно выполняется несколько запросов, каждый из которых получает не более 100 элементов.
// Количество одновременно работающих горутин можно контроллировать через goroutineLimit (по умолчанию 1)
func (s searchInstance[T]) AllAtOnce(ctx context.Context, goroutineLimit int) ([]T, error) {

	count, err := s.Count(ctx)
	if err != nil {
		return nil, err
	}

	numberOfCycles := 1 + count/100

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if goroutineLimit < 1 {
		goroutineLimit = 1
	}

	eg, _ := errgroup.WithContext(ctx)
	eg.SetLimit(goroutineLimit)

	all := make([]T, 0, count)
	mu := sync.Mutex{}
	for i := 0; i < numberOfCycles; i++ {
		i := i
		eg.Go(func() error {
			items, _, err := s.app.find(ctx, filter{
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
				return nil
			}
			mu.Lock()
			defer mu.Unlock()
			all = append(all, items...)
			return nil
		})
	}

	err = eg.Wait()
	if err == nil || errors.Is(err, ErrNoMoreItems) {
		return all, nil
	}
	return nil, err
}

// First получает один элмент по переданному фильтру
func (s searchInstance[T]) First(ctx context.Context) (T, error) {

	var t T
	items, _, err := s.app.find(ctx, filter{
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

// Count возвращает общее кол-во элементов по переданному фильтру.
// Аналог COUNT в SQL.
func (s searchInstance[T]) Count(ctx context.Context) (int, error) {

	_, count, err := s.app.find(ctx, filter{
		From:         s.from,
		Size:         0,
		Active:       !s.includeDeleted,
		SearchFilter: s.search,
	})
	if err != nil {
		return 0, err
	}

	return count, nil

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

// Size позволяет регулировать максимальное кол-во элментов,
// которые будут возвращены при поиске (но не более 100, по умолчанию 10).
// Аналог LIMIT в SQL
func (s searchInstance[T]) Size(size int) searchInstance[T] {
	if size < 0 {
		size = 10
	}
	s.size = size
	return s
}

// From позволяет регулировать с какого по счету элемента будет выполнен поиск.
// Аналог OFFSET в SQL
func (s searchInstance[T]) From(from int) searchInstance[T] {
	if from < 0 {
		from = 0
	}
	s.from = from
	return s
}

// IncludeDeleted добавляет выборку удаленные элменты (__deletedAt != null)
func (s searchInstance[T]) IncludeDeleted() searchInstance[T] {
	s.includeDeleted = true
	return s
}
