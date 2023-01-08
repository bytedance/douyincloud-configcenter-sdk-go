package base

import (
	"time"
)

const (
	DefaultUpdateInterval = 1 * time.Minute
	DefaultTimeout        = 5 * time.Second
)

type clientConfig struct {
	updateInterval time.Duration
	timeout        time.Duration
}

func NewClientConfig() *clientConfig {
	c := &clientConfig{
		updateInterval: DefaultUpdateInterval,
		timeout:        DefaultTimeout,
	}
	return c
}

func (c *clientConfig) SetUpdateInterval(t time.Duration) {
	if t < 10*time.Second {
		t = 10 * time.Second
	}
	c.updateInterval = t
}

func (c *clientConfig) SetTimeout(t time.Duration) {
	c.timeout = t
}

func (c *clientConfig) GetUpdateInterval() time.Duration {
	return c.updateInterval
}

func (c *clientConfig) GetTimeout() time.Duration {
	return c.timeout
}
