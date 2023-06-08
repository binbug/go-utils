package httputils

import (
	"errors"
	"github.com/binbug/go-utils/jsonutils"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Get[T any](URL string, opts ...Option) (httpResult HttpResult[T]) {
	cfg := initConfig(opts...)
	client := &http.Client{
		Timeout: cfg.timeout,
	}

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return HttpResult[T]{err: err}
	}

	if cfg.header != nil {
		req.Header = cfg.header
	}

	resp, err := client.Do(req)
	if err != nil {
		return HttpResult[T]{err: err}
	}

	return processResponse[T](resp, cfg)
}

func PostForm[T any](URL string, data url.Values, opts ...Option) (httpResult HttpResult[T]) {
	return Post[T](URL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()), opts...)
}

func PostJSON[T any](URL string, o interface{}, opts ...Option) (httpResult HttpResult[T]) {
	return Post[T](URL, "application/json", strings.NewReader(jsonutils.ToJSON(o)), opts...)
}

func Post[T any](URL string, contentType string, body io.Reader, opts ...Option) (httpResult HttpResult[T]) {
	cfg := initConfig(opts...)

	client := &http.Client{
		Timeout: cfg.timeout,
	}

	req, err := http.NewRequest("POST", URL, body)

	if err != nil {
		return HttpResult[T]{err: err}
	}

	if cfg.header != nil {
		req.Header = cfg.header
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := client.Do(req)
	if err != nil {
		return HttpResult[T]{err: err}
	}

	return processResponse[T](resp, cfg)
}

func initConfig(opts ...Option) *Config {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func processResponse[T any](resp *http.Response, cfg *Config) (httpResult HttpResult[T]) {

	httpResult = HttpResult[T]{
		statusCode: resp.StatusCode,
	}

	switch any(httpResult.o).(type) {
	case *http.Response:
		httpResult.o = any(resp).(T)
		return httpResult
	case http.Response:
		httpResult.o = any(*resp).(T)
		return httpResult
	case io.ReadCloser:
		httpResult.o = any(resp.Body).(T)
		return httpResult
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		httpResult.err = err
		return httpResult
	}

	if cfg.rawBody {
		httpResult.rawBody = data
	}

	httpResult.statusCode = resp.StatusCode
	httpResult.header = resp.Header
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		httpResult.err = errors.New(string(data))
		return httpResult
	}

	o, err := jsonutils.FromBytes[T](data)
	if err != nil {
		httpResult.err = err
		return httpResult
	}

	httpResult.o = o
	return httpResult
}
