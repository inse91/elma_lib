package e365_gateway

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUrlGeneration(t *testing.T) {
	appSettings := AppSettings{
		Stand: NewStand(StandSettings{
			Host:  "https://elma.ru",
			Port:  "8080",
			Token: "",
		}),
		Namespace: "ns1",
		Code:      "app1",
	}
	const wantUrl = "https://elma.ru:8080/pub/v1/app/ns1/app1"
	require.Equal(t, wantUrl, appSettings.toAppUrl())
}
