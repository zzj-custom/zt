package ncm

import (
	"encoding/json"
	"io"
	"net/http"
	"zt/backend/pkg/request"
)

type ImageRequest struct {
	request.Request
}

func (i ImageRequest) Validate() error {
	return nil
}

func (i ImageRequest) GetUrl() string {
	return i.Request.Url
}

func (i ImageRequest) GetBody() io.Reader {
	return nil
}

func (i ImageRequest) GetMethod() string {
	return http.MethodGet
}

func (i ImageRequest) GetHeaders() map[string]string {
	return nil
}

func (i ImageRequest) GetTimeout() int {
	return 30
}

type ImageResponse struct {
	Data []byte `json:"data"`
}

func (i ImageResponse) IsOk() bool {
	return true
}

func (i ImageResponse) GetMessage() string {
	return ""
}

func (i ImageResponse) Decode(v any) error {
	return json.Unmarshal(i.Data, v)
}
