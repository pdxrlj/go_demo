package main

import (
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

var rooms = sync.Map{}

type Message struct {
	From      any `json:"from"`
	To        any `json:"to"`
	Candidate any `json:"candidate"`
	Offer     any `json:"offer"`
	Answer    any `json:"answer"`
	OffLine   any `json:"offline"`
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		user := ws.Request().URL.Query().Get("user")
		if _, ok := rooms.Load(user); !ok {
			rooms.Store(user, ws)
			rooms.Range(func(key, value any) bool {
				if user == key {
					return true
				}
				_ = websocket.JSON.Send(ws, &Message{
					From: key,
					To:   "all",
				})

				return true
			})
		}
		for {
			var message Message
			err := websocket.JSON.Receive(ws, &message)
			if err != nil {
				rooms.Delete(user)
				rooms.Range(func(key, value any) bool {
					if conn, ok := value.(*websocket.Conn); ok {
						err := websocket.JSON.Send(conn, &Message{
							From:    "system",
							To:      key,
							OffLine: user,
						})
						if err != nil {
							log.Println("write json:", err)
							return false
						}
					}
					return true
				})
				break
			}

			if message.To == "all" {
				rooms.Range(func(key, value any) bool {

					if conn, ok := value.(*websocket.Conn); ok {
						err := websocket.JSON.Send(conn, message)
						if err != nil {
							log.Println("write json:", err)
							return false
						}
					}
					return true
				})
				continue
			}

			if to, ok := rooms.Load(message.To); ok {
				if conn, ok := to.(*websocket.Conn); ok {
					err := websocket.JSON.Send(conn, message)
					if err != nil {
						log.Println("write json:", err)
						break
					}
				}
			}
		}

	}))

	mux.Handle("/", http.FileServer(http.Dir("./html")))

	log.Fatal(http.ListenAndServeTLS(":3000", "127.pem", "127_key.pem", mux))
}
