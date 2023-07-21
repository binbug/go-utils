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

	deserializeCode map[int]struct{}

	checkRedirect func(req *http.Request, via []*http.Request) error
}

type Option func(*Config)

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.timeout = timeout
	}
}

func WithCheckRedirect(fn func(req *http.Request, via []*http.Request) error) Option {
	return func(c *Config) {
		c.checkRedirect = fn
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

func DeserializeCode(codes ...int) Option {
	return func(c *Config) {
		if c.deserializeCode == nil {
			c.deserializeCode = make(map[int]struct{})
		}

		for _, code := range codes {
			c.deserializeCode[code] = struct{}{}
		}
	}
}

func withReqInterceptor(fn func(req *http.Request)) Option {
	return func(c *Config) {
		c.reqInterceptor = fn
	}
}
