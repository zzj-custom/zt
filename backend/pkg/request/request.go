package request

import (
	"io"
)

const (
	Timeout = 30
)

type RequestInterface interface {
	// Validate 验证数据
	Validate() error

	// GetUrl 获取请求地址
	GetUrl() string

	// GetBody 获取请求结构体
	GetBody() io.Reader

	// GetMethod 获取请求方法
	GetMethod() string

	// GetHeaders 获取请求头
	GetHeaders() map[string]string

	// GetTimeout 获取超时时间
	GetTimeout() int
}

type Request struct {
	Url     string    `json:"url"`
	Body    io.Reader `json:"body"`
	Method  string    `json:"method"`
	Timeout int       `json:"timeout"`
}

type ResponseInterface interface {
	// IsOk 是否成功
	IsOk() bool

	// GetMessage 获取状态码
	GetMessage() string

	// Decode 解析数据
	Decode(v any) error
}
