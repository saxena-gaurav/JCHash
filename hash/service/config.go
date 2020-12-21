package service

import (
	"fmt"
	"time"
)

const (
	defaultHashDelay  = 5
	defaultListenPort = 8080
)

type Config struct {
	ListenPort int
	HashDelay  int
}

func (c *Config) hashDelaySeconds() time.Duration {
	if c.HashDelay == 0 {
		c.HashDelay = defaultHashDelay
	}
	return time.Duration(c.HashDelay) * time.Second
}

func (c *Config) listenAddr() string {
	if c.ListenPort == 0 {
		c.ListenPort = defaultListenPort
	}
	return fmt.Sprintf(":%d", c.ListenPort)
}
