package openapi

type HttpResp struct {
	Code int                   `json:"code"`
	Msg  string                `json:"msg"`
	Data GetConfigListResponse `json:"data"`
}

type GetConfigListResponse struct {
	Kvs     []ItemInfo `json:"kvs"`
	Version string     `json:"version"`
}

type ItemInfo struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  int64  `json:"type"`
}

type GetConfigListReqBody struct {
	Version string `json:"version"`
}
