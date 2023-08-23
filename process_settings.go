package e365_gateway

const pubV1ApiInstance = "pub/v1/bpm/instance/"
const pubV1ApiBpm = "pub/v1/bpm/template"

type ProcessSettings struct {
	Stand     Stand
	Namespace string
	Code      string
}
