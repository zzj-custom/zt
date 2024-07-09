package request

import (
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type Config struct {
	Host string `json:"host" toml:"host"`
}

type Client struct {
	cfg *Config
}

func (r *Client) Do(httpReq RequestInterface, httpResp ResponseInterface) error {
	if err := httpReq.Validate(); err != nil {
		return errors.Wrap(err, "invalid request parameter")
	}
	url := r.buildApiUrl(httpReq)

	bodyByte, err := r.getBody(httpReq)
	if err != nil {
		return errors.Wrap(err, "get request body failed")
	}

	req, err := http.NewRequest(httpReq.GetMethod(), url, httpReq.GetBody())
	if err != nil {
		return errors.Wrap(err, "create request failed")
	}
	cli := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

	st := time.Now().UnixMilli()
	response, err := cli.Do(req)
	ed := time.Now().UnixMilli()
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	respBytes, err := io.ReadAll(response.Body)

	s := slog.With(
		slog.String("response", string(respBytes)),
		slog.String("url", url),
		slog.String("method", httpReq.GetMethod()),
		slog.String("cost", fmt.Sprintf("%dms", ed-st)),
		slog.String("body", string(bodyByte)),
	)

	if tid, ok := httpReq.GetHeaders()["X-TRACE-ID"]; ok {
		s = s.With(slog.String("trace_id", tid))
	}

	defer func() {
		s.Info("请求耗时")
		_ = response.Body.Close()
	}()
	if err != nil {
		return errors.Wrap(err, "read response failed")
	}

	if err := httpResp.Decode(respBytes); err != nil {
		return errors.Wrap(err, "decode response failed")
	}

	if !httpResp.IsOk() {
		return errors.New(httpResp.GetMessage())
	}
	return nil
}

func (r *Client) getBody(httpReq RequestInterface) (body []byte, err error) {
	if httpReq.GetBody() == nil {
		return nil, nil
	}
	return io.ReadAll(httpReq.GetBody())
}

func (r *Client) buildApiUrl(req RequestInterface) string {
	return fmt.Sprintf("%s%s", r.cfg.Host, req.GetUrl())
}

var (
	client     *Client
	clientOnce sync.Once
)

func NewClient(cfg *Config) *Client {
	clientOnce.Do(func() {
		client = &Client{
			cfg: cfg,
		}
	})
	return client
}
