package rooms

import (
	"fmt"
	"net"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/json-iterator/go"
)

type Message struct {
	// 消息类型
	Type string `json:"type"`
	// 消息内容
	Content string `json:"content"`
	// 消息发送者
	From string `json:"from"`
	// 消息接收者
	To string `json:"to"`
}

type Rooms struct {
	rooms map[string]*Users
	lock  sync.RWMutex
}

func NewDefaultRooms() *Rooms {
	return &Rooms{
		rooms: make(map[string]*Users),
		lock:  sync.RWMutex{},
	}
}

func (r *Rooms) AddToRooms(room, uid string, c net.Conn) *Users {
	r.lock.Lock()
	defer r.lock.Unlock()
	if _, ok := r.rooms[room]; !ok {
		r.rooms[room] = NewUsers()
	}
	r.rooms[room].AddToUsers(uid, c)
	return r.rooms[room]
}

func (r *Rooms) DisconnectUser(room, uid string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	users := r.rooms[room]
	users.RemoveUser(uid)
}

func (r *Rooms) RemoveRooms(room string, uid string) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	rooms := r.rooms[room]
	if rooms == nil {
		return
	}
	rooms.RemoveUser(uid)
}

func (r *Rooms) Broadcast(rooms string, msg []byte) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	errsArr := make([]error, 0)
	users := r.rooms[rooms]

	users.Each(func(c net.Conn, uid string) error {
		if err := wsutil.WriteServerMessage(c, ws.OpText, msg); err != nil {
			return err
		}
		return nil
	}, nil)

	if len(errsArr) > 0 {
		panic(errsArr)
	}
}

func (r *Rooms) SendTo(rooms, uid, msg string) error {
	r.lock.RLock()
	defer r.lock.RUnlock()
	bytesContent, err := jsoniter.MarshalToString(&Message{
		Type:    "sendTo",
		Content: msg,
		From:    rooms,
		To:      uid,
	})
	if err != nil {
		panic(err)
	}
	users := r.rooms[rooms]
	return users.Emit(uid, bytesContent)
}

func (r *Rooms) Exists(room, uid string) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	users := r.rooms[room]
	if users == nil {
		return false
	}
	return users.Exists(uid)
}

func (r *Rooms) GetRooms(room string) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	users := r.rooms[room]
	if users == nil {
		fmt.Printf("room %s not exists\n", room)
		return
	}
	var uids []string
	users.Each(func(c net.Conn, uid string) error {
		uids = append(uids, uid)
		return nil
	}, nil)

	bytesContent, err := jsoniter.MarshalToString(&uids)
	if err != nil {
		panic(err)
	}
	sendMessage := &SendMessage{
		Type:    "lists",
		Content: bytesContent,
		From:    "server",
		To:      "all",
	}

	bytes, _ := sendMessage.Bytes()

	r.Broadcast(room, bytes)
}
