package e365_gateway

import (
	"fmt"
)

const pubV1ApiApp = "pub/v1/app"

//type AppCommon struct {
//	URL   string
//	Token string
//}

type AppSettings struct {
	//AppCommon
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
