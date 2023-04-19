package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gobwas/ws"

	"websocket/rooms"
)

const (
	defaultRoom = "default"
)

var defaultRoomsManager = rooms.NewDefaultRooms()

func main() {
	fmt.Printf("websocket runing in 8080\n")

	go func() {
		for {
			defaultRoomsManager.GetRooms(defaultRoom)
			time.Sleep(5 * time.Second)
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			panic(err)
		}

		// 获取请求的参数uid
		uid := r.URL.Query().Get("uid")
		if uid == "" {
			panic("uid is empty")
		}
		extra := r.URL.Query().Get("extra")
		if extra == "" {
			panic("extra is empty")
		}
		if exists := defaultRoomsManager.Exists(defaultRoom, uid); exists {
			return
		}

		fmt.Printf("new user conn:%v\n", uid)
		// add to rooms
		users := defaultRoomsManager.AddToRooms(defaultRoom, uid, extra, conn)

		go func(uid string) {
			defer func(uid string) {
				defaultRoomsManager.DisconnectUser(defaultRoom, uid)
			}(uid)

			for {
				if err := users.Read(uid, func(content *rooms.ReceiveMessage) error {
					sendMessage := &rooms.SendMessage{
						From:    content.From,
						To:      content.To,
						Type:    content.Type,
						Content: content.Content,
					}
					fmt.Printf("read send message: %+v from:%v to:%v\n", sendMessage.Type, sendMessage.From, sendMessage.To)
					// 发送消息
					users.Write(sendMessage)
					return nil
				}); err != nil {
					fmt.Printf("read error: %+v\n", err)
					panic(err)
				}
			}
		}(uid)
	})))
}
