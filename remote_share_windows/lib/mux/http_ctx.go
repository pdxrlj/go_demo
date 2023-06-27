package mux

import (
	"net/http"

	"remote_share_windows/config"
)

type Ctx struct {
	w http.ResponseWriter
	r *http.Request
	c *config.AppConfig
}

func newCtx() *Ctx {
	return &Ctx{}
}

func (c *Ctx) AddHttpResponseWriter(w http.ResponseWriter) *Ctx {
	c.w = w
	return c
}

func (c *Ctx) AddHttpRequest(r *http.Request) *Ctx {
	c.r = r
	return c
}

func (c *Ctx) LoadAppConfig(config *config.AppConfig) *Ctx {
	c.c = config
	return c
}

func (c *Ctx) GetHttpResponseWriter() http.ResponseWriter {
	return c.w
}

func (c *Ctx) GetHttpRequest() *http.Request {
	return c.r
}

func (c *Ctx) GetAppConfig() *config.AppConfig {
	return c.c
}
