package httputils

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/binbug/go-utils/jsonutils"
)

// Get function is used to send an HTTP GET request.
// T represents the type of the response body.
// URL represents the URL of the request.
// opts represents the optional configuration options.
func Get[T any](URL string, opts ...Option) (httpResult HttpResult[T]) {
	return execute[T](URL, http.MethodGet, nil, opts...)
}

// PostForm function is used to send an HTTP POST request with form data.
// T represents the type of the response body.
// URL represents the URL of the request.
// data represents the form data to be sent.
// opts represents the optional configuration options.
func PostForm[T any](URL string, data url.Values, opts ...Option) (httpResult HttpResult[T]) {
	return Post[T](URL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()), opts...)
}

// PostJSON function is used to send an HTTP POST request with JSON data.
// T represents the type of the response body.
// URL represents the URL of the request.
// o represents the JSON data to be sent.
// opts represents the optional configuration options.
func PostJSON[T any](URL string, o interface{}, opts ...Option) (httpResult HttpResult[T]) {
	return Post[T](URL, "application/json", strings.NewReader(jsonutils.ToJSON(o)), opts...)
}

// Post function is used to send an HTTP POST request.
// T represents the type of the response body.
// URL represents the URL of the request.
// contentType represents the content type of the request body.
// body represents the request body.
// opts represents the optional configuration options.
func Post[T any](URL string, contentType string, body io.Reader, opts ...Option) (httpResult HttpResult[T]) {
	opts = append(opts, withReqInterceptor(func(req *http.Request) {
		req.Header.Set("Content-Type", contentType)
	}))

	return execute[T](URL, http.MethodPost, body, opts...)
}

// Delete is a function that sends an HTTP DELETE request to the specified URL.
//
// Parameters:
//
//	URL string - the URL to send the DELETE request to.
//	opts ...Option - optional options for the HTTP request.
//
// Return:
//
//	HttpResult[T] - the result of the HTTP request.
func Delete[T any](URL string, opts ...Option) HttpResult[T] {
	return execute[T](URL, http.MethodDelete, nil, opts...)
}

// execute is a function that sends an HTTP request based on the provided URL, method, body, and options.
//
// URL: the URL to send the HTTP request to.
// method: the HTTP method to use for the request.
// body: the request body.
// opts: optional configurations for the request.
// Returns an HttpResult containing the result of the HTTP request.
func execute[T any](URL, method string, body io.Reader, opts ...Option) HttpResult[T] {
	cfg := initConfig(opts...)
	client := &http.Client{
		Timeout:       cfg.timeout,
		CheckRedirect: cfg.checkRedirect,
	}

	req, err := http.NewRequest(method, URL, body)
	if err != nil {
		return HttpResult[T]{err: err}
	}

	if cfg.header != nil {
		req.Header = cfg.header
	}

	for _, interceptor := range cfg.reqInterceptors {
		interceptor(req)
	}

	resp, err := client.Do(req)
	if err != nil {
		return HttpResult[T]{err: err}
	}

	return processResponse[T](resp, cfg)

}

// initConfig function is used to initialize the configuration for an HTTP request.
// opts represents the optional configuration options.
// Returns a pointer to the initialized Config struct.
func initConfig(opts ...Option) *Config {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// processResponse function is used to process an HTTP response.
// T represents the type of the response body.
// resp represents the HTTP response to be processed.
// cfg represents the configuration options for the HTTP request.
// Returns an HttpResult struct containing the processed response.
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
	case struct{}:
		httpResult.o = any(struct{}{}).(T)
		return httpResult
	case *struct{}:
		httpResult.o = any(&struct{}{}).(T)
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
	if len(cfg.deserializeCode) > 0 {
		_, ok := cfg.deserializeCode[resp.StatusCode]
		if !ok {
			httpResult.err = errors.New(string(data))
			return httpResult
		}
	} else if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
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
