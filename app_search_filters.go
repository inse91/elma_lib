package e365_gateway

import (
	"encoding/json"
	"strings"
	"time"
)

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
		i++
	}
	sb.WriteString("}}")
	return []byte(sb.String()), nil
}

type appNumberFilter map[string]float64

var AppNumberFilter = appNumberFilter{
	"min": -(1 << 63),
	"max": 1 << 63,
}

func (mf appNumberFilter) To(to float64) appNumberFilter {
	mf["max"] = to
	return mf
}
func (mf appNumberFilter) From(from float64) appNumberFilter {
	mf["min"] = from
	return mf
}

type appDateFilter map[string]string

var AppDateFilter = appDateFilter{
	"min": "1970-01-01",
	"max": "3000-01-01",
}

func (adf appDateFilter) From(from time.Time) appDateFilter {
	adf["min"] = from.Format(time.DateOnly)
	return adf
}

func (adf appDateFilter) To(to time.Time) appDateFilter {
	adf["max"] = to.Format(time.DateOnly)
	return adf
}

// AppDateFilter позволяет делать фильтр по дате
//func AppDateFilter(min, max string) map[string]string {
//	return map[string]string{
//		"min": min,
//		"max": max,
//	}
//}

// AppNumberFilter фильтрует по числовому полю - OLD
//func AppNumberFilter(min, max float64) map[string]float64 {
//	return map[string]float64{
//		"min": min,
//		"max": max,
//	}
//}

func AppCategoryFilter(code string) string {
	return code
}

func AppApplicationFilter(id string) [1]string {
	return [1]string{id}
}
