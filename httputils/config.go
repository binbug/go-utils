package httputils

import (
	"net/http"
	"time"
)

type Config struct {
	timeout time.Duration

	header http.Header

	rawBody bool

	reqInterceptor func(req *http.Request)
}

type Option func(*Config)

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.timeout = timeout
	}
}

func WithHeader(header http.Header) Option {
	return func(c *Config) {
		c.header = header
	}
}

func WithRawBody() Option {
	return func(c *Config) {
		c.rawBody = true
	}
}

func withReqInterceptor(fn func(req *http.Request)) Option {
	return func(c *Config) {
		c.reqInterceptor = fn
	}
}
