package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"ws_message/client"
)

type MessageType string

const (
	Broadcast MessageType = "broadcast"
	Group     MessageType = "group"
	Private   MessageType = "private"
)

// Message 消息结构体
// Private 消息发送者和接收者必须一致
// Group   接受者为群组id
// Broadcast 接受者全群广播
type Message struct {
	// 消息类型
	Type MessageType `json:"type"`
	// 消息内容
	Content any `json:"content"`
	// 消息发送者
	Sender string `json:"sender"`
	// 消息接收者 room:uuid
	Receiver string `json:"receiver"`
	// 消息发送时间
	Time string `json:"time"`
	// 消息的发送uuid
	MessageUUID string `json:"message_uuid"`
}

func (m *Message) GetType() MessageType {
	return m.Type
}

func (m *Message) GetContent() any {
	return m.Content
}

func (m *Message) GetSender() string {
	return m.Sender
}

type MessageHandler func(data []byte) error

type MessageHandlerI interface {
	Apply(*Message)
}

type BroadcastMessageHandler func(message *Message)

func (b BroadcastMessageHandler) Apply(message *Message) {
	b(message)
}

func WithBroadcastMessageHandler(server *Server) MessageHandlerI {
	return BroadcastMessageHandler(func(message *Message) {
		if message.GetType() != Broadcast {
			return
		}
		for room, clients := range server.clients {
			clients.Range(func(key, value interface{}) bool {
				if c, ok := value.(*client.Client); ok {
					if err := c.SendMessage(message, client.JsonClientType); err != nil {
						fmt.Printf("send broadcast message err: %s\n", err.Error())
						server.RemoveClient(room, key.(string))
						c.Close()
					}
				}
				return true
			})
		}
	})
}

type GroupMessageHandler func(message *Message)

func (g GroupMessageHandler) Apply(message *Message) {
	g(message)
}

func WithGroupMessageHandler(server *Server) MessageHandlerI {
	return GroupMessageHandler(func(message *Message) {
		if message.GetType() != Group {
			return
		}
		if clients, ok := server.clients[message.Receiver]; ok {
			fmt.Printf("send group message: %+v\n", message)
			clients.Range(func(key, value interface{}) bool {
				if c, ok := value.(*client.Client); ok {
					if err := c.SendMessage(message, client.JsonClientType); err != nil {
						fmt.Printf("send group message err: %s\n", err.Error())
						server.RemoveClient(message.Receiver, key.(string))
						c.Close()
					}
				}
				return true
			})
		}
	})
}

type PrivateMessageHandler func(message *Message)

func (p PrivateMessageHandler) Apply(message *Message) {
	p(message)
}

func WithPrivateMessageHandler(server *Server) MessageHandlerI {
	return PrivateMessageHandler(func(message *Message) {
		if message.GetType() != Private {
			return
		}
		// room:uuid
		if room, uuid, ok := strings.Cut(message.Receiver, ":"); ok {
			fmt.Printf("send private room:%v uuid:%v message: %+v\n", room, uuid, message)
			if clients, ok := server.clients[room]; ok {
				clients.Range(func(key, value interface{}) bool {
					if c, ok := value.(*client.Client); ok {
						if c.UUID() == uuid {
							if err := c.SendMessage(message, client.JsonClientType); err != nil {
								fmt.Printf("send private message err: %s\n", err.Error())
								server.RemoveClient(message.Receiver, key.(string))
								c.Close()
							}
							// stop range
							return false
						}
					}
					return true
				})
			}
		}

	})
}

func (m *Message) HandlerReceiver(handlers ...MessageHandlerI) any {
	for _, handler := range handlers {
		m.MessageUUID = uuid.New().String()
		m.Time = time.Now().Format(time.DateTime)
		handler.Apply(m)
	}

	return m.Content
}
