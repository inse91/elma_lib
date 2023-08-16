package e365_gateway

type StatusInfo struct {
	StatusItems []StatusItem `json:"statusItems"`
	GroupItems  []GroupItem  `json:"groupItems"`
}

type StatusItem struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	GroupId string `json:"groupId"`
}

type GroupItem struct {
	Id   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
