package request

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
	// "go.elastic.co/apm"
	// "go.elastic.co/apm/module/apmhttp"
)

// Options request options
type Options struct {
	method  string
	data    []byte
	headers map[string]string
	timeout int
	ssl     bool
	ctx     context.Context
}

// Option request option
type Option func(*Options)

func initOptions(options ...Option) *Options {
	opts := &Options{
		method:  http.MethodGet,
		timeout: 15,
		headers: map[string]string{"Content-Type": "application/json; charset=UTF-8"},
		ssl:     true,
		ctx:     context.Background(),
	}
	for _, option := range options {
		option(opts)
	}
	return opts
}

// WithOptions accepts the whole options config.
func WithOptions(options Options) Option {
	return func(opts *Options) {
		*opts = options
	}
}

// WithMethod set request method.
func WithMethod(method string) Option {
	return func(opts *Options) {
		opts.method = method
	}
}

// WithData set request data.
func WithData(data []byte) Option {
	return func(opts *Options) {
		opts.data = data
	}
}

// WithHeader set request header.
func WithHeader(header map[string]string) Option {
	return func(opts *Options) {
		opts.headers = header
	}
}

// WithTimeout set request timeout.
func WithTimeout(timeout int) Option {
	return func(opts *Options) {
		opts.timeout = timeout
	}
}

// WithSSL set request skip ssl verify.
func WithSSL(ssl bool) Option {
	return func(opts *Options) {
		opts.ssl = ssl
	}
}

// WithContext set context.
func WithContext(ctx context.Context) Option {
	return func(opts *Options) {
		opts.ctx = ctx
	}
}

// DoRequest exec https? request and return []byte
func DoRequest(url string, options ...Option) (code int, respBuf []byte, respHeader map[string][]string, err error) {
	// exec the undercourse request
	resp, err := DoRequestUndercourse(url, options...)
	if err != nil {
		// 错误
		return -1, nil, nil, errors.New("response failure")
	}
	respBuf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, resp.Header, errors.New("read response failure")
	}
	defer resp.Body.Close()

	return resp.StatusCode, respBuf, resp.Header, nil
}

// DoRequestUndercourse exec https? request and return response
func DoRequestUndercourse(url string, options ...Option) (resp *http.Response, err error) {
	opts := initOptions(options...)
	// span, ctx := apm.StartSpan(opts.ctx, "dorequest", "custom")
	// defer span.End()

	var req *http.Request
	switch opts.method {
	case http.MethodPost, http.MethodPut:
		req, err = http.NewRequest(opts.method, url, bytes.NewBuffer(opts.data))
	default:
		req, err = http.NewRequest(opts.method, url, nil)
	}
	if err != nil {
		return nil, errors.New("build request failure")
	}

	for key, val := range opts.headers {
		req.Header.Set(key, val)
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: opts.ssl,
			},
			DisableKeepAlives: true,
		},
		Timeout: time.Duration(opts.timeout) * time.Second,
		// CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// 	return http.ErrUseLastResponse
		// },
	}
	// client = apmhttp.WrapClient(client)
	// resp, err = client.Do(req.WithContext(ctx))
	resp, err = client.Do(req.WithContext(opts.ctx))
	return resp, err
}