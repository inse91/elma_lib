package e365_gateway

type respCommon struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type itemResponse[T interface{}] struct {
	respCommon
	Item T `json:"item"`
}

type createItemRequest[T interface{}] struct {
	Context T `json:"context"`
}

type setStatusRequest struct {
	Status statusCode `json:"status"`
}

type statusCode struct {
	Code string `json:"code"`
}

type getStatusResponse struct {
	respCommon
	StatusInfo
}

type appListResponse[T interface{}] struct {
	respCommon
	Result appListResult[T] `json:"result"`
}

type appListResult[T interface{}] struct {
	Result []T `json:"result"`
	Total  int `json:"total"`
}

type runProcRequest[T interface{}] struct {
	Context T `json:"context"`
}

type runProcResponse[T interface{}] struct {
	respCommon
	Context T `json:"context"`
}

type getProcInstanceResponse[T interface{}] struct {
	respCommon
	Context T `json:"data"`
}

type getFileLinkResp struct {
	respCommon
	Link string `json:"Link"`
}

type fileResponse struct {
	respCommon
	File File `json:"file"`
}

type dirInfoResponse struct {
	respCommon
	Directory DirectoryInfo `json:"directory"`
}
