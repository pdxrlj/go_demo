package handler

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"golang.org/x/sync/errgroup"

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

	request.AssetPort = query.Get("asset_port")
	request.AssetHost = query.Get("asset_host")
	request.AssetProtocol = query.Get("asset_protocol")
	request.AssetUser = query.Get("asset_user")
	request.AssetPassword = query.Get("asset_password")
	request.GuacamoleAddr = query.Get("guacamole_addr")

	request.ScreenWidth = cast.ToInt(query.Get("screen_width"))
	request.ScreenHeight = cast.ToInt(query.Get("screen_height"))
	request.ScreenDpi = cast.ToInt(query.Get("screen_dpi"))

	protocol := ctx.GetHttpRequest().Header.Get(http.CanonicalHeaderKey("Sec-Websocket-Protocol"))
	conn, err := upgrade.Upgrade(ctx.GetHttpResponseWriter(), ctx.GetHttpRequest(), http.Header{
		"Sec-Websocket-Protocol": []string{protocol},
	})
	if err != nil {
		fmt.Printf("升级ws失败: %v\n", err)
		return err
	}

	// close ws
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("关闭ws失败: %v\n", err)
		}
	}()

	// 2. 连接远程桌面
	uid := ""
	tunnel, err := lib.NewGuacamoleTunnel(&lib.GuacamoleTunnel{
		GuacamoleAddr: request.GuacamoleAddr,
		Protocol:      request.AssetProtocol,
		Host:          request.AssetHost,
		Port:          request.AssetPort,
		User:          request.AssetUser,
		Password:      request.AssetPassword,
		Uuid:          uid,
		W:             request.ScreenWidth,
		H:             request.ScreenHeight,
		Dpi:           request.ScreenDpi,
	})
	if err != nil {
		fmt.Printf("连接远程桌面失败: %v\n", err)
		return err
	}

	// close tunnel
	defer func() {
		err := tunnel.Close()
		if err != nil {
			fmt.Printf("关闭tunnel失败: %v\n", err)
		}
	}()

	ioCopy(conn, tunnel)

	return nil
}

var InternalOpcodeIns = []byte(fmt.Sprint(len(""), ".", ""))

func ioCopy(conn *websocket.Conn, tunnel *lib.SimpleTunnel) {
	writer := tunnel.AcquireWriter()
	reader := tunnel.AcquireReader()

	defer func() {
		tunnel.ReleaseWriter()
		tunnel.ReleaseReader()
	}()

	eg := errgroup.Group{}
	eg.Go(func() error {
		buf := bytes.NewBuffer(make([]byte, 0, lib.MaxGuacamoleMessage*2))
		for {
			read, err := reader.Read()
			if err != nil {
				fmt.Printf("reader.Read err: %v\n", err)
				return err
			}

			if bytes.HasPrefix(read, InternalOpcodeIns) {
				// messages starting with the InternalDataOpcode are never sent to the websocket
				continue
			}

			_, err = buf.Write(read)

			err = conn.WriteMessage(websocket.TextMessage, buf.Bytes())
			if err != nil {
				fmt.Printf("conn.WriteMessage err: %v\n", err)
				return err
			}
			buf.Reset()
		}
	})

	eg.Go(func() error {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("conn.ReadMessage err: %v\n", err)
				return err
			}

			if bytes.HasPrefix(data, InternalOpcodeIns) {
				// messages starting with the InternalDataOpcode are never sent to guacd
				continue
			}
			if _, err = writer.Write(data); err != nil {
				fmt.Printf("ReadMessage after writer.Write err: %v\n", err)
				return err
			}
		}
	})

	if err := eg.Wait(); err != nil {
		fmt.Printf("ioCopy err: %v\n", err)
		return
	}
}
