package lib

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
)

type Stream struct {
	Conn         net.Conn
	parseStart   int
	ConnectionID string
	buffer       []byte
	reset        []byte
	timeout      time.Duration
}

func NewStream(conn net.Conn, timeout time.Duration) *Stream {
	buffer := make([]byte, 0, MaxGuacamoleMessage*3)
	return &Stream{
		Conn:    conn,
		timeout: timeout,
		buffer:  buffer,
		reset:   buffer[:cap(buffer)],
	}
}

func (s *Stream) Handshake(config *GuacamoleConfig) (err error) {
	selectArg := config.ConnectionID
	if len(selectArg) == 0 {
		selectArg = config.Protocol
	}
	_, err = s.Write(NewInstruction("select", selectArg).Bytes())
	if err != nil {
		return err
	}

	read, err := s.Read()
	if err != nil {
		panic(err)
		return err
	}
	fmt.Printf("read: %s\n", read)

	return nil
}

func (s *Stream) Write(data []byte) (int, error) {
	err := s.Conn.SetWriteDeadline(time.Now().Add(s.timeout))
	if err != nil {
		return 0, errors.WithStack(err)
	}
	fmt.Printf("")
	return s.Conn.Write(data)
}

func (s *Stream) Read() ([]byte, error) {
	err := s.Conn.SetReadDeadline(time.Now().Add(s.timeout))
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(s.Conn)

	for {
		readByte, err := reader.ReadBytes(';')
		if err != nil {
			fmt.Printf("=========end err: %s\n", err)
			return nil, err
		}

		fmt.Printf("readByte: %s\n", string(readByte))
	}
}

func (s *Stream) AssertOpcode(opcode string) (*Instruction, error) {
	Parse(s)
	return nil, nil
}
