package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/json-iterator/go"
)

var rooms = sync.Map{}

type Message struct {
	From      any `json:"from"`
	To        any `json:"to"`
	Candidate any `json:"candidate"`
	Offer     any `json:"offer"`
	Answer    any `json:"answer"`
}

func main() {
	app := fiber.New()

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		user := c.Query("user")
		if _, ok := rooms.Load(user); !ok {
			rooms.Store(user, c)
		}

		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			log.Printf("recv: %s", msg)

			if len(msg) <= 0 {
				continue
			}
			var message Message
			err = jsoniter.Unmarshal(msg, &message)
			fmt.Printf("message=%v\n", message)
			if err == nil {
				if message.To == "all" {
					rooms.Range(func(key, value any) bool {
						fmt.Printf("brocast to all key=%v, value=%v\n", key, value)
						err := value.(*websocket.Conn).WriteMessage(mt, msg)
						if err != nil {
							log.Println("write json:", err)
							return false
						}
						return true
					})
					continue
				}

				if to, ok := rooms.Load(message.To); ok {
					err := to.(*websocket.Conn).WriteMessage(mt, msg)
					if err != nil {
						log.Println("write json:", err)
						break
					}
				}
			} else {
				err = c.WriteMessage(mt, msg)
				if err != nil {
					log.Println("write:", err)
					break
				}
			}

		}
	}))

	log.Fatal(app.Listen(":3000"))
}
