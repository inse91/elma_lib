package e365_gateway

import (
	"fmt"
	"net/http"
)

type stand struct {
	host string
	port string
	//_token string
	h http.Header
}

var testDefaultStandSettings = StandConfig{
	Host:  "https://q3bamvpkvrulg.elma365.ru",
	Port:  "",
	Token: "33ef3e66-c1cd-4d99-9a77-ddc4af2893cf",
}

type StandConfig struct {
	Host  string
	Port  string
	Token string
}

func (s stand) url() string {
	if s.port == "" {
		return s.host
	}
	return fmt.Sprintf("%s:%s", s.host, s.port)
}

func (s stand) header() http.Header {
	return s.h
}

type Stand interface {
	url() string
	//token() string
	header() http.Header
}

func NewStand(settings StandConfig) Stand {
	return stand{
		host: settings.Host,
		port: settings.Port,
		h: func() http.Header {
			h := http.Header{}
			h.Set("Content-type", "application/json")
			h.Set("Authorization", "Bearer "+settings.Token)
			return h
		}(),
	}
}
