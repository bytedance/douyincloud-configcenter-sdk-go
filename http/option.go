package http

import (
	"net/http"
	"time"
)

type Option func(options *Options)

type Options struct {
	Timeout   time.Duration
	Url       string
	Header    http.Header
	Retry     uint64
	Transport *http.Transport
}

func DefaultOptions() *Options {
	header := http.Header{}
	header.Set("Content-Type", "application/json")

	return &Options{
		Timeout: 5 * time.Second,
		Header:  header,
		Transport: &http.Transport{
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     10 * time.Second,
		},
		Retry: 3,
	}
}
