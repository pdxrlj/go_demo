package handler

import (
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"

	"remote_share_windows/lib"
	"remote_share_windows/lib/mux"
)

func ApiWsGuacamole(ctx *mux.Ctx) error {
	websocketReadBufferSize := lib.MaxGuacamoleMessage
	websocketWriteBufferSize := lib.MaxGuacamoleMessage * 2
	upgrade := websocket.Upgrader{
		ReadBufferSize:  websocketReadBufferSize,
		WriteBufferSize: websocketWriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	var request RequestGuacamole
	// 1. 获取请求的http query参数
	query := ctx.GetHttpRequest().URL.Query()
	err := sonic.Unmarshal([]byte(query.Encode()), &request)
	if err != nil {
		return err
	}

	protocol := ctx.GetHttpRequest().Header.Get(http.CanonicalHeaderKey("Sec-Websocket-Protocol"))
	conn, err := upgrade.Upgrade(ctx.GetHttpResponseWriter(), ctx.GetHttpRequest(), http.Header{
		"Sec-Websocket-Protocol": []string{protocol},
	})
	if err != nil {
		return err
	}

	// close ws
	defer func() {
		_ = conn.Close()
	}()

	// 2. 连接远程桌面
	//uid := ""

	return nil
}
