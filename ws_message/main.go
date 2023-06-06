package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"

	"ws_message/client"
	"ws_message/server"
)

func main() {
	mux := http.NewServeMux()

	chatWs := server.NewServer()
	defer chatWs.Stop()
	// websocket
	mux.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		fmt.Printf("client connect: %s\n", ws.Request().RemoteAddr)
		// websocket 连接
		uuid := ws.Request().URL.Query().Get("uuid")
		if uuid == "" {
			fmt.Printf("client err uuid: %s\n", uuid)
			return
		}
		chatWs.AddClient("", client.NewClient(
			client.SetClientConn(ws),
			client.SetClientUUID(uuid),
			client.SetClientMaxHeartbeatTimes(3),
		))
	}))

	mux.HandleFunc("/ws/users/lists", func(writer http.ResponseWriter, request *http.Request) {
		rooms := request.URL.Query().Get("rooms")
		users := chatWs.GetRoomUsers(rooms)
		_, _ = writer.Write(users)
	})

	log.Fatal(http.ListenAndServe(":8181", mux))

}
