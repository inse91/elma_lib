package e365_gateway

const uuid4Len = 36

type Elma struct {
	url   string
	token string
}

func New(url, token string) Elma {
	return Elma{
		url:   url,
		token: token,
	}
}
