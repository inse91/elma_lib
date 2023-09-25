package e365_gateway

import (
	"encoding/json"
	"strings"
	"time"
)

// filter общий набор фильтров
type filter struct {
	From   int  `json:"from"`
	Size   int  `json:"size"`
	Active bool `json:"active"`
	SearchFilter
}

// SearchFilter - набор фильтров для выполнения поиска
type SearchFilter struct {
	Fields          Fields           `json:"filter"`
	IDs             []string         `json:"ids,omitempty"`
	SortExpressions []SortExpression `json:"sortExpressions,omitempty"`
	AtStatus        []string         `json:"statusCode,omitempty"`
	StatusGroupId   string           `json:"statusGroupId,omitempty"`
}

type SortExpression struct {
	Ascending bool   `json:"ascending"`
	Field     string `json:"field"`
}

// Fields - обертка для фильтров по полям приложений.
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
		i++
	}
	sb.WriteString("}}")
	return []byte(sb.String()), nil
}

type fielder interface {
	Category(code string) string
	App(id string) [1]string
	DateTime() appDateFilter
	Number() appNumberFilter
}

// Field - обертка для поиска по полям типов "Число", "Дата", "Категория" и "Приложение"
var Field fielder = fieldIface{}

type fieldIface struct{}

// Category фильтр для полей типа "Категория"
func (f fieldIface) Category(code string) string {
	return code
}

// App фильтр для полей типа "Приложение"
func (f fieldIface) App(id string) [1]string {
	return [1]string{id}
}

// Date фильтр для полей типа "Дата"
func (f fieldIface) DateTime() appDateFilter {
	return appDateFilter{
		"min": "1970-01-01",
		"max": "3000-01-01",
	}
}

// Number фильтр для полей типа "Число"
func (f fieldIface) Number() appNumberFilter {
	return appNumberFilter{
		"min": -(1 << 63),
		"max": 1 << 63,
	}
}

type appNumberFilter map[string]float64

// To позволяет задать верхнюю границу для полей типа "Число" (опционально)
func (mf appNumberFilter) To(value float64) appNumberFilter {
	mf["max"] = value
	return mf
}

// From позволяет задать нижнюю границу для полей типа "Число" (опционально)
func (mf appNumberFilter) From(value float64) appNumberFilter {
	mf["min"] = value
	return mf
}

// Equal позволяет задать равенство для полей типа "Число" (опционально)
func (mf appNumberFilter) Equal(value float64) appNumberFilter {
	mf["min"] = value
	mf["max"] = value
	return mf
}

type appDateFilter map[string]string

// From позволяет задать верхнюю границу для полей типа "Дата" (опционально)
func (adf appDateFilter) From(date time.Time) appDateFilter {
	adf["min"] = date.UTC().Format(time.RFC3339)
	return adf
}

// To позволяет задать верхнюю границу для полей типа "Дата" (опционально)
func (adf appDateFilter) To(date time.Time) appDateFilter {
	adf["max"] = date.UTC().Format(time.RFC3339)
	return adf
}

// EqualDate позволяет задать равенство для полей типа "Дата" (опционально) с точностью до одного дня
func (adf appDateFilter) EqualDate(date time.Time) appDateFilter {
	oneDay := time.Hour * 24
	start := date.Truncate(oneDay)
	end := start.Add(oneDay)
	return adf.From(start).To(end)
}
