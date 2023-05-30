package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"golang.org/x/sync/errgroup"
)

func main() {
	engine := gin.Default()
	wsServer := &WsServer{}
	engine.GET("/ws", wsServer.Handler)

	log.Fatal(engine.Run(":8089"))

}

type WsServer struct {
}

func (w *WsServer) Handler(ctx *gin.Context) {
	w.handler(ctx)
}

func (w *WsServer) handler(ctx *gin.Context) {
	ws := websocket.Server{}
	ws.Handler = func(conn *websocket.Conn) {
		serverConn, err := w.dialServer(ctx.Request)
		if err != nil {
			panic(err)
		}

		eg := errgroup.Group{}
		eg.SetLimit(2)
		eg.Go(func() error {
			// client -> server
			_, err := io.Copy(serverConn, conn)
			if err != nil {
				return err
			}
			serverConn.Close()
			return nil
		})
		eg.Go(func() error {
			// server -> client
			_, err := io.Copy(conn, serverConn)
			if err != nil {
				return err
			}
			conn.Close()
			return nil
		})

		err = eg.Wait()
		if err != nil {
			panic(err)
		}

	}

	ws.ServeHTTP(ctx.Writer, ctx.Request)
}

func (w *WsServer) dialServer(req *http.Request) (*websocket.Conn, error) {
	origin := req.Header.Get("Origin")
	if len(origin) == 0 {
		panic("no origin")
	}
	fmt.Printf("origin is %s \n", origin)

	config, err := websocket.NewConfig("ws://127.0.0.1:3000/ws", origin)
	if err != nil {
		panic(err)
	}

	config.Header = req.Header.Clone()

	config.Header.Del("Sec-WebSocket-Extensions")

	config.Header.Del("Origin")

	const xForwardedFor = "X-Forwarded-For"
	xff := req.Header.Get(xForwardedFor)
	if host, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if xff == "" {
			config.Header.Set(xForwardedFor, host)
		} else {
			config.Header.Set(xForwardedFor, xff+", "+host)
		}
	}

	const xForwardedHost = "X-Forwarded-Host"
	xfh := req.Header.Get(xForwardedHost)
	if xfh == "" && req.Host != "" {
		config.Header.Set(xForwardedHost, req.Host)
	}

	const xForwardedProto = "X-Forwarded-Proto"
	config.Header.Set(xForwardedProto, "http")
	if req.TLS != nil {
		config.Header.Set(xForwardedProto, "https")
	}

	return websocket.DialConfig(config)
}
