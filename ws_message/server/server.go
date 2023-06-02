package server

import (
	"fmt"
	"sync"
	"time"

	"ws_message/client"

	"github.com/bytedance/sonic"
)

const (
	defaultRoom = "default"
)

type Server struct {
	// sync.Map uuid -> client
	clients    map[string]*sync.Map
	clientNums int
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) Stop() *Server {
	for _, clients := range server.clients {
		clients.Range(func(key, value interface{}) bool {
			// 断开连接
			if c, ok := value.(*client.Client); ok {
				c.Close()
			}
			return true
		})
	}
	return server
}

func (server *Server) AddClient(room string, c *client.Client) *Server {
	if room == "" {
		room = defaultRoom
	}

	if server.clients == nil {
		server.clients = make(map[string]*sync.Map)
	}

	if _, ok := server.clients[room]; !ok {
		server.clients[room] = &sync.Map{}
	}

	server.clients[room].Store(c.UUID(), c)

	// 监听用户心跳
	go func() {
		ticker := time.NewTicker(c.GetHeartbeatInterval())
		for {
			select {
			case <-ticker.C:
				ticker.Stop()
				if err := c.Heartbeat(); err != nil {
					fmt.Printf("heartbeat error: %+v\n", err)
					server.RemoveClient(room, c.UUID())
					c.Close()
					return
				}
				//fmt.Printf("heartbeat success: %+v\n", c.UUID())
				ticker.Reset(c.GetHeartbeatInterval())
			}
		}
	}()

	// 监听用户消息
	if err := c.HandlerMessage(func(data []byte) error {
		var message Message
		if err := sonic.Unmarshal(data, &message); err != nil {
			fmt.Printf("unmarshal message error: %+v\n", err)
		} else {
			message.HandlerReceiver(
				WithBroadcastMessageHandler(server),
				WithGroupMessageHandler(server),
				WithPrivateMessageHandler(server),
			)
		}

		return nil
	}); err != nil {
		fmt.Printf("handler message error: %+v\n", err)
		server.RemoveClient(room, c.UUID())
		c.Close()
	}

	return server
}

func (server *Server) RemoveClient(room string, uuid string) *Server {
	if room == "" {
		room = defaultRoom
	}

	if clients, ok := server.clients[room]; ok {
		clients.Delete(uuid)
	}
	return server
}
