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

func NewStand(host, port, token string) Stand {
	return stand{
		host: host,
		port: port,
		h: func() http.Header {
			h := http.Header{}
			h.Set("Content-type", "application/json")
			h.Set("Authorization", "Bearer "+token)
			return h
		}(),
	}
}
