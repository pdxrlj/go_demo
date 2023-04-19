package rooms

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
)

type user struct {
	c     net.Conn
	tick  *time.Ticker
	extra string
}

type Users struct {
	users map[string]*user
	lock  sync.RWMutex
}

func NewUsers() *Users {
	return &Users{
		users: make(map[string]*user),
		lock:  sync.RWMutex{},
	}
}

func (u *Users) PingCheck(uid string) error {
	u.lock.RLock()
	defer u.lock.RUnlock()

	if u.users[uid] == nil {
		return nil
	}
	ping := &SendMessage{
		Type:    "ping",
		Content: "ping",
		From:    "server",
		To:      uid,
	}
	bytes, err := jsoniter.Marshal(ping)
	if err != nil {
		return errors.WithStack(err)
	}
	return wsutil.WriteServerMessage(u.users[uid].c, ws.OpText, bytes)
}

func (u *Users) AddToUsers(uid, extra string, c net.Conn) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.users[uid] = &user{
		c:     c,
		tick:  time.NewTicker(15 * time.Second),
		extra: extra,
	}
	// 发送心跳包
	go func() {
		for {
			if u.users[uid] != nil {
				select {
				case <-u.users[uid].tick.C:
					if err := u.PingCheck(uid); err != nil {
						fmt.Printf("ping uid:%v error: %v\n", uid, err)
						u.RemoveUser(uid)
						return
					}
				}
			}
		}
	}()
}

func (u *Users) RemoveUser(uid string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	// 停止心跳包
	if u.users[uid] == nil {
		return
	}

	u.users[uid].tick.Stop()
	_ = u.users[uid].c.Close()
	delete(u.users, uid)
}

func (u *Users) Emit(uid string, msg string) error {
	u.lock.RLock()
	defer u.lock.RUnlock()
	bytesContent, err := jsoniter.Marshal(&Message{
		Type:    "text",
		Content: msg,
		From:    uid,
		To:      uid,
	})
	if err != nil {
		panic(err)
	}
	if err = wsutil.WriteServerMessage(u.users[uid].c, ws.OpText, bytesContent); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (u *Users) Each(fn func(c net.Conn, uid string) error, errHandler func(err error)) {
	u.lock.RLock()
	defer u.lock.RUnlock()
	for id, conn := range u.users {
		if err := fn(conn.c, id); err != nil {
			if errHandler == nil {
				fmt.Printf("each conn error: %v\n", err)
				continue
			}
			errHandler(err)
		}
	}
}

type ReadHandler[T any] func(content T) error

type NextReadHandler[T any] func(next ReadHandler[T]) ReadHandler[T]

func HandlerReadContent[T any](content T, handler ReadHandler[T], next ...NextReadHandler[T]) error {
	for i := len(next) - 1; i >= 0; i-- {
		handler = next[i](handler)
	}

	return handler(content)
}

type ReceiveMessage struct {
	// 消息类型
	Type string `json:"type"`
	// 消息内容
	Content any `json:"content"`
	// 消息发送者
	From string `json:"from"`
	// 消息接收者
	To string `json:"to"`
}

type SendMessage struct {
	// 消息类型
	Type string `json:"type"`
	// 消息内容
	Content any `json:"content"`
	// 消息发送者
	From string `json:"from"`
	// 消息接收者
	To string `json:"to"`
}

func (s *SendMessage) Bytes() ([]byte, error) {
	return jsoniter.Marshal(s)
}

func DefaultConvertReceiverMessage[T []byte](content T) (*ReceiveMessage, error) {
	if content == nil {
		return nil, errors.New("content is nil")
	}
	var receiverMessage ReceiveMessage
	if err := jsoniter.Unmarshal(content, &receiverMessage); err != nil {
		return nil, err
	}
	return &receiverMessage, nil
}

func (u *Users) Read(uid string, handler ReadHandler[*ReceiveMessage]) error {
	u.lock.RLock()
	var message *ReceiveMessage
	if u.users[uid] == nil {
		u.lock.RUnlock()
		return nil
	}
	u.lock.RUnlock()

	data, opCode, err := wsutil.ReadClientData(u.users[uid].c)
	if err != nil {
		fmt.Printf("read err type:%T\n", err)
		u.RemoveUser(uid)
		return nil
	}
	if err = u.checkUserMessage(opCode, data); err != nil {
		if strings.Contains(err.Error(), "close") {
			u.RemoveUser(uid)
			return nil
		}
		return err
	}

	return HandlerReadContent(message, func(content *ReceiveMessage) error {
		return handler(content)
	}, func(next ReadHandler[*ReceiveMessage]) ReadHandler[*ReceiveMessage] {
		return func(content *ReceiveMessage) error {
			message, err = DefaultConvertReceiverMessage(data)
			if err != nil {
				return errors.WithStack(err)
			}
			//fmt.Printf("read message: %+v\n", message)
			return next(message)
		}
	}, func(next ReadHandler[*ReceiveMessage]) ReadHandler[*ReceiveMessage] {
		return func(content *ReceiveMessage) error {
			// 过滤消息类型
			// candidate 消息交换
			if content.Type == "candidate" {
				return next(content)
			}

			// offer 消息交换
			if content.Type == "offer" {
				return next(content)
			}

			return nil
		}
	})
}

func (u *Users) Write(message *SendMessage) {
	bytesContent, err := jsoniter.Marshal(message)
	if err != nil {
		panic(err)
	}
	if u.users[message.To] == nil {
		fmt.Printf("write to user %v is not exist\n", message.To)
		return
	}

	if err = wsutil.WriteServerMessage(u.users[message.To].c, ws.OpText, bytesContent); err != nil {
		fmt.Printf("write to user %v error: %v\n", message.To, err)
	}
}

func (u *Users) checkUserMessage(opCode ws.OpCode, data []byte) error {
	if opCode == ws.OpClose {
		return errors.New("close")
	}

	if opCode == ws.OpPing {
		if err := wsutil.WriteServerMessage(u.users["1"].c, ws.OpPong, data); err != nil {
			return errors.WithStack(err)
		}
	}

	if opCode == ws.OpPong {
		if err := wsutil.WriteServerMessage(u.users["1"].c, ws.OpPing, data); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (u *Users) Exists(uid string) bool {
	if _, exists := u.users[uid]; exists {
		return true
	}
	return false
}
