package http

import (
	"context"
	"net/http"
	"os"
)

const (
	HttpMethodPost = "POST"
	DefaultUrl     = "http://config-center.dycloud.run:8699/config/get_config_list"
	EnvUrl         = "DYC_CONFIG_CENTER_URL"
)

type Client struct {
	opts *Options
	c    http.Client
}

func NewClient(options ...Option) *Client {
	url := GetEnvDefault(EnvUrl, DefaultUrl)
	opts := DefaultOptions()
	opts.Url = url
	for _, o := range options {
		o(opts)
	}
	return &Client{
		opts: opts,
		c: http.Client{
			Transport: opts.Transport,
		},
	}
}

func (client *Client) GetOptions() *Options {
	return client.opts
}

type APIInfo struct {
	Method string
	Body   string
	Header http.Header
}

// HttpPostRaw 发起JSON的post请求
func (client *Client) HttpPostRaw(body string, headers http.Header) ([]byte, http.Header, error) {
	return client.CtxHttpPostRaw(context.Background(), body, headers)
}

func (client *Client) CtxHttpPostRaw(ctx context.Context, body string, headers http.Header) ([]byte, http.Header, error) {
	apiInfo := &APIInfo{
		Method: HttpMethodPost,
		Body:   body,
		Header: headers,
	}
	return client.request(ctx, apiInfo)
}

func GetEnvDefault(key, defVal string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	return val
}
