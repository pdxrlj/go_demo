package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"golang.org/x/net/websocket"
	"golang.org/x/sync/errgroup"
)

func TestConn(t *testing.T) {
	rand.NewSource(time.Now().UnixNano())
	uuid := strconv.Itoa(rand.Intn(1000000))

	dial, err := websocket.Dial("ws://localhost:8181/ws?uuid="+uuid, "", "http://127.0.0.1")

	if err != nil {
		t.Error(err)
		return
	}

	randMessage := func() ([]byte, error) {
		m := struct {
			Type string `json:"type"`
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
		}{
			Type:    "broadcast",
			Content: RandomString(10000),
			Sender:  uuid,
		}

		return sonic.Marshal(&m)
	}
	eg := errgroup.Group{}
	eg.SetLimit(10000)
	eg.Go(func() error {
		readMessage(dial)
		return nil
	})

	for i := 0; i < 100000; i++ {
		eg.Go(func() error {
			message, err := randMessage()
			if err != nil {
				return err
			}
			_, err = dial.Write(message)

			return err
		})
	}
	if err := eg.Wait(); err != nil {
		t.Error(err)
		return
	}
}

func readMessage(conn *websocket.Conn) {
	for {
		message := make([]byte, 1024)
		read, err := conn.Read(message)
		if err != nil {
			return
		}
		if read < len(message) {
			fmt.Printf("read end")
			return
		}

		fmt.Printf("read message: %s\n", message)
	}

}

func RandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}
