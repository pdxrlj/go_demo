package server

import (
	"net/http"

	"remote_share_windows/config"
	"remote_share_windows/lib/mux"
	"remote_share_windows/server/handler"
)

func NewRouter() *mux.Http {
	r := mux.NewHttp(mux.WithAppConfig(&config.AppConfig{}))
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
			w.Header().Set("content-type", "application/json")             //返回数据格式是json
			next.ServeHTTP(w, r)
		})
	})
	r.Get("/ws", handler.ApiWsGuacamole)

	return r
}
