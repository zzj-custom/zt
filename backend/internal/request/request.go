package request

import (
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	cookiemonster "github.com/MercuryEngineering/CookieMonster"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"
)

var (
	// 默认请求头
	defaultHeaders = map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Charset":  "UTF-8,*;q=0.5",
		"Accept-Encoding": "gzip,deflate,sdch",
		"Accept-Language": "en-US,en;q=0.8",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36",
	}
)

type ReqOptions struct {
	RetryTimes int
	Cookie     string
	UserAgent  string
	Refer      string
	Silent     bool
}

type Request struct {
	options *ReqOptions
}

// NewRequest 初始化一个请求
func NewRequest(opt *ReqOptions) *Request {
	return &Request{
		options: &ReqOptions{
			RetryTimes: opt.RetryTimes,
			Cookie:     opt.Cookie,
			UserAgent:  opt.UserAgent,
			Refer:      opt.Refer,
			Silent:     opt.Silent,
		},
	}
}

func DefaultRequest() *Request {
	return NewRequest(&ReqOptions{
		RetryTimes: 3,
	})
}

// Do base request
func (r *Request) Do(method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DisableCompression:  true,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Minute,
		Jar:       jar,
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for k, v := range defaultHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if _, ok := headers["Referer"]; !ok {
		req.Header.Set("Referer", url)
	}
	if r.options.Cookie != "" {
		// 添加cookie在请求里面
		cookies, _ := cookiemonster.ParseString(r.options.Cookie)
		if len(cookies) > 0 {
			for _, c := range cookies {
				req.AddCookie(c)
			}
		} else {
			// cookie is not Netscape HTTP format, set it directly
			// a=b; c=d
			req.Header.Set("Cookie", r.options.Cookie)
		}
	}

	if r.options.UserAgent != "" {
		req.Header.Set("User-Agent", r.options.UserAgent)
	}

	if r.options.Refer != "" {
		req.Header.Set("Referer", r.options.Refer)
	}

	var (
		res *http.Response
	)
	for i := 0; ; i++ {
		res, err = client.Do(req)
		if err == nil && res.StatusCode < 400 {
			break
		} else if i+1 >= r.options.RetryTimes {
			if err != nil {
				err = errors.Errorf("request error: %v", err)
			} else {
				err = errors.Errorf("%s request error: HTTP %d", url, res.StatusCode)
			}
			return nil, errors.WithStack(err)
		}
		time.Sleep(1 * time.Second)
	}
	return res, nil
}

// Get request
func (r *Request) Get(url, refer string, headers map[string]string) (string, error) {
	body, err := r.GetByte(url, refer, headers)
	return string(body), err
}

// GetByte get request
func (r *Request) GetByte(url, refer string, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	if refer != "" {
		headers["Referer"] = refer
	}
	res, err := r.Do(http.MethodGet, url, nil, headers)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() { _ = res.Body.Close() }()

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(res.Body)
	case "deflate":
		reader = flate.NewReader(res.Body)
	default:
		reader = res.Body
	}
	defer func() { _ = reader.Close() }()

	body, err := io.ReadAll(reader)
	if err != nil && err != io.EOF {
		return nil, errors.WithStack(err)
	}
	return body, nil
}

// Headers 返回地址的请求头
func (r *Request) Headers(url, refer string) (http.Header, error) {
	headers := map[string]string{
		"Referer": refer,
	}
	res, err := r.Do(http.MethodGet, url, nil, headers)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() { _ = res.Body.Close() }()
	return res.Header, nil
}

// Size 获取地址的文本内容长度
func (r *Request) Size(url, refer string) (int64, error) {
	h, err := r.Headers(url, refer)
	if err != nil {
		return 0, errors.Wrapf(err, "get size of %s", url)
	}
	s := h.Get("Content-Length")
	if s == "" {
		return 0, errors.New("Content-Length is not exists")
	}
	size, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "parse size of %s", url)
	}
	return size, nil
}

// ContentType 获取地址的Content-Type
func (r *Request) ContentType(url, refer string) (string, error) {
	h, err := r.Headers(url, refer)
	if err != nil {
		return "", errors.Wrapf(err, "get heaer of %s", url)
	}
	s := h.Get("Content-Type")
	if s == "" {
		return "", errors.New("Content-Type is not exists")
	}
	return strings.Split(s, ";")[0], nil
}
