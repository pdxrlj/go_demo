package mux

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"remote_share_windows/config"
)

type Http struct {
	*chi.Mux
	c *config.AppConfig
}

type Option func(*Http)

type Handler func(ctx *Ctx) error

func WithDefaultHttpMux() Option {
	return func(h *Http) {
		h.Mux = chi.NewRouter()
	}
}

func WithAppConfig(config *config.AppConfig) Option {
	return func(h *Http) {
		h.c = config
	}
}

func (h *Http) Get(path string, handler Handler) {
	h.Mux.Get(path, func(writer http.ResponseWriter, request *http.Request) {
		err := handler(newCtx().
			AddHttpResponseWriter(writer).
			AddHttpRequest(request).
			LoadAppConfig(h.c))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (h *Http) Post(path string, handler Handler) {
	h.Mux.Post(path, func(writer http.ResponseWriter, request *http.Request) {
		err := handler(newCtx().
			AddHttpResponseWriter(writer).
			AddHttpRequest(request).
			LoadAppConfig(h.c))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (h *Http) Put(path string, handler Handler) {
	h.Mux.Put(path, func(writer http.ResponseWriter, request *http.Request) {
		err := handler(newCtx().
			AddHttpResponseWriter(writer).
			AddHttpRequest(request).
			LoadAppConfig(h.c))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (h *Http) Delete(path string, handler Handler) {
	h.Mux.Delete(path, func(writer http.ResponseWriter, request *http.Request) {
		err := handler(newCtx().
			AddHttpResponseWriter(writer).
			AddHttpRequest(request).
			LoadAppConfig(h.c))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func NewHttp(opts ...Option) *Http {
	h := &Http{
		Mux: chi.NewMux(),
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}
