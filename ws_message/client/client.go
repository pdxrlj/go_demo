package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
)

const (
	defaultMaxHeartbeatTimes = 3
	headerHeartbeatInterval  = time.Second * 3
)

type (
	MessageHandler func(message []byte) error
)

type Client struct {
	// user identity
	uuid string
	// websocket connection
	conn *websocket.Conn
	// heartbeat interval if heartbeat than  maxHeartbeatTimes remove the client

	maxHeartbeatTimes int
	//heartbeat interval
	heartbeatInterval time.Duration
	// last heartbeat time
	lasTimeHeartbeat time.Time

	closeChannel   chan struct{}
	MessageHandler MessageHandler
	once           sync.Once
}

type Option func(*Client)

func newClient() *Client {
	return &Client{
		maxHeartbeatTimes: defaultMaxHeartbeatTimes,
		closeChannel:      make(chan struct{}),
		heartbeatInterval: headerHeartbeatInterval,
	}
}

// NewClient create a new client
func NewClient(options ...Option) *Client {
	client := newClient()
	for _, option := range options {
		option(client)
	}
	return client
}

func (c *Client) GetHeartbeatInterval() time.Duration {
	return c.heartbeatInterval
}

func (c *Client) MaxHeartbeatTimes() int {
	return c.maxHeartbeatTimes
}

func (c *Client) UUID() string {
	return c.uuid
}

// SetClientUUID set Client uuid
func SetClientUUID(uuid string) Option {
	return func(client *Client) {
		client.uuid = uuid
	}
}

// SetClientConn set Client conn
func SetClientConn(conn *websocket.Conn) Option {
	return func(client *Client) {
		client.conn = conn
	}
}

// SetClientMaxHeartbeatTimes set Client maxHeartbeatTimes
func SetClientMaxHeartbeatTimes(maxHeartbeatTimes int) Option {
	return func(client *Client) {
		client.maxHeartbeatTimes = maxHeartbeatTimes
	}
}

// SetClientLasTimeHeartbeat set Client lasTimeHeartbeat
func SetClientLasTimeHeartbeat(lasTimeHeartbeat time.Time) Option {
	return func(client *Client) {
		client.lasTimeHeartbeat = lasTimeHeartbeat
	}
}

// HandlerMessage read message from client
// if read error than close client and server should remove client
func (c *Client) HandlerMessage(handler MessageHandler) error {
	for {
		select {
		case <-c.closeChannel:
			err := c.conn.Close()
			if err != nil {
				return err
			}
		default:
			var msg []byte
			err := websocket.Message.Receive(c.conn, &msg)
			if err != nil {
				if errors.Is(err, websocket.ErrBadClosingStatus) {
					fmt.Printf("client %s have close\n", c.uuid)
					_ = c.conn.Close()
					return websocket.ErrBadClosingStatus
				}
				return errors.WithStack(err)
			}
			fmt.Printf("1 client %s receive message: %s\n", c.uuid, msg)
			if err := handler(msg); err != nil {
				fmt.Printf("2 client %s receive message: %s\n", c.uuid, msg)
				return err
			}
		}
	}
}

// SendMessage send message to client
// if send error and ErrBadClosingStatus error
// client and server should remove client
func (c *Client) SendMessage(message any, messageType ...ClientType) error {
	mType := bytesClientType
	if len(messageType) > 0 {
		mType = messageType[0]
	}
	switch mType {
	case JsonClientType:
		return websocket.JSON.Send(c.conn, message)
	default:
		return websocket.Message.Send(c.conn, message)
	}
}

// Heartbeat send heartbeat to client
func (c *Client) Heartbeat() error {
	var err error
	if err = c.SendMessage([]byte("heartbeat")); err != nil {
		fmt.Printf("send heartbeat client:%v retry:%d interval:%+v \n", c.uuid, c.MaxHeartbeatTimes(), c.GetHeartbeatInterval())
		for i := 0; i < c.MaxHeartbeatTimes(); i++ {
			if err = c.SendMessage("heartbeat"); err != nil {
				fmt.Printf("send heartbeat error: %+v\n", err)
				continue
			}
		}
	}
	return err
}

func (c *Client) Close() {
	fmt.Printf("client %+v close\n", c.uuid)
	c.once.Do(func() {
		close(c.closeChannel)
	})
}
