package e365_gateway

import (
	"fmt"
)

const pubV1ApiApp = "pub/v1/app"
const pubV1ApiBpm = "pub/v1/bpm/template"

//type Common struct {
//	URL   string
//	Token string
//}

type AppSettings struct {
	//Common
	Stand     Stand
	Namespace string
	Code      string
}

func (s AppSettings) toAppUrl() string {
	return fmt.Sprintf("%s/%s/%s/%s", s.Stand.url(), pubV1ApiApp, s.Namespace, s.Code)
}

func (s AppSettings) toBpmUrl() string {
	return fmt.Sprintf("%s/%s/%s/%s", s.Stand.url(), pubV1ApiBpm, s.Namespace, s.Code)
}
