package e365_gateway

import (
	"fmt"
)

type Settings struct {
	Stand     Stand
	Namespace string
	Code      string
}

func (s Settings) toAppUrl() string {
	if s.Stand == nil {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s/%s", s.Stand.url(), pubV1ApiApp, s.Namespace, s.Code)
}

func (s Settings) toBpmUrl() string {
	if s.Stand == nil {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s/%s", s.Stand.url(), pubV1ApiBpm, s.Namespace, s.Code)
}
